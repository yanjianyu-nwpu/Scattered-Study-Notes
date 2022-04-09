# 笔记

## Sever 端

### 基本结构体

```'
// 提前标注 空指针错误，不能直接使用error 因为 反射reflect 是空interface
var typeOfError = reflect.TypeOf((*error)(nil)).Elem()
// 每一个方法的信息
type methodType struct {
	// 这个锁主要是用于维护numcalls，并发安全
	sync.Mutex // protects counters
	method     reflect.Method
	ArgType    reflect.Type
	ReplyType  reflect.Type
	// 召唤的总数
	numCalls   uint
}
// service 的情况，是所有的method的集合。是结构体
type service struct {
	name   string                 // name of service
	rcvr   reflect.Value          // receiver of methods for the service
	typ    reflect.Type           // type of the receiver
	method map[string]*methodType // registered methods
}
// Request is a header written before every RPC call. It is used internally
// Request 结构体是一个头部 在rpc call之前写入的，内部使用。现在是存储起来用于debug比如分细 network 分细
// but documented here as an aid to debugging, such as when analyzing
// network traffic.
type Request struct {
	ServiceMethod string   // format: "Service.Method"
	Seq           uint64   // sequence number chosen by client
	next          *Request // for free list in Server
}
// 同上
type Response struct {
	ServiceMethod string    // echoes that of the Request
	Seq           uint64    // echoes that of the request
	Error         string    // error, if any.
	next          *Response // for free list in Server
}
// Server represents an RPC Server.
// server结构体待变一个rpc server
type Server struct {
	// 用于 登记名字 和 service 结构体的  
	serviceMap sync.Map   // map[string]*service
	// 用于保护 freeReq
	reqLock    sync.Mutex // protects freeReq
	freeReq    *Request
	respLock   sync.Mutex // protects freeResp
	freeResp   *Response
}
// DefaultServer 是一个单例。且是上述的Server的实例
var DefaultServer = NewServer()
```

后面学一学 go/token这个包

### 注册函数

```Go
// Register publishes in the server the set of methods of the
// receiver value that satisfy the following conditions:
// Register 发布
//	- exported method of exported type
//  - 导出的方法和类型
//	- two arguments, both of exported type
// - 两个参数，都是导出的类型
//	- the second argument is a pointer
// - 第二个参数是指针
//	- one return value, of type error
// - 一个返回值。是error
// It returns an error if the receiver is not an exported type or has
// no suitable methods. It also logs the error using package log.
// 返回一个错误，如果receiver 不是一个导出类型，或者没有合适的method 。也会打log
// The client accesses each method using a string of the form "Type.Method",
// 每个 客户端 进入
// where Type is the receiver's concrete type.
func (server *Server) Register(rcvr interface{}) error {
	return server.register(rcvr, "", false)
}
// 这里 是真实的使用
func (server *Server) register(rcvr interface{}, name string, useName bool) error {		
	// 这个是结构体 级别的
	s := new(service)
	// 赋值 type
	s.typ = reflect.TypeOf(rcvr)
	// 赋值值
	s.rcvr = reflect.ValueOf(rcvr)
	// reflect.Indirect返回指向对象的类型，然后记录名字，
	// 然后这样用reflect.Indiredct 获取真实的type name
	sname := reflect.Indirect(s.rcvr).Type().Name()

	if useName {
		sname = name
	}
	if sname == "" {
		s := "rpc.Register: no service name for type " + s.typ.String()
		log.Print(s)
		return errors.New(s)
	}
	// 如果是导出的
	if !token.IsExported(sname) && !useName {
		s := "rpc.Register: type " + sname + " is not exported"
		log.Print(s)
		return errors.New(s)
	}
	s.name = sname

	// Install the methods
	s.method = suitableMethods(s.typ, true)

	if len(s.method) == 0 {
		str := ""

		// To help the user, see if a pointer receiver would work.
		method := suitableMethods(reflect.PtrTo(s.typ), false)
		if len(method) != 0 {
			str = "rpc.Register: type " + sname + " has no exported methods of suitable type (hint: pass a pointer to value of that type)"
		} else {
			str = "rpc.Register: type " + sname + " has no exported methods of suitable type"
		}
		log.Print(str)
		return errors.New(str)
	}

	if _, dup := server.serviceMap.LoadOrStore(sname, s); dup {
		return errors.New("rpc: service already defined: " + sname)
	}
	return nil
}

然后这里的 suitableMthods 看函数是否符合效果，一个返回error 两个参数，一个指针一个返回值


```

然后后续就是编码解码和 发送

```
// getResponse 获取一个Response
func (server *Server) getResponse() *Response {
    // 先上锁
	server.respLock.Lock()
	// 这里感觉就是获取需要的resp，然后发出去
	resp := server.freeResp
	if resp == nil {
		resp = new(Response)
	} else {
		server.freeResp = resp.next
		*resp = Response{}
	}
	server.respLock.Unlock()
	return resp
}

// 这就是 sendResponse 这是发送执行后的结果
func (server *Server) sendResponse(sending *sync.Mutex, req *Request, reply interface{}, codec ServerCodec, errmsg string) {
	resp := server.getResponse()
	// Encode the response header
	resp.ServiceMethod = req.ServiceMethod
	if errmsg != "" {
		resp.Error = errmsg
		reply = invalidRequest
	}
	resp.Seq = req.Seq
	sending.Lock()
	// 然后write codec 结果
	err := codec.WriteResponse(resp, reply)
	if debugLog && err != nil {
		log.Println("rpc: writing response:", err)
	}
	sending.Unlock()
	server.freeResponse(resp)
}
```

这里是真正之后发送出去

```

// call 调用
func (s *service) call(server *Server, sending *sync.Mutex, wg *sync.WaitGroup, mtype *methodType, req *Request, argv, replyv reflect.Value, codec ServerCodec) {
	if wg != nil {
		defer wg.Done()
	}
	mtype.Lock()
	mtype.numCalls++
	mtype.Unlock()
	// 然后这里找到对应函数
	function := mtype.method.Func
	// 调用函数，并返回一个新的返回体
	// Invoke the method, providing a new value for the reply.
	returnValues := function.Call([]reflect.Value{s.rcvr, argv, replyv})
	// The return value for the method is an error.
	errInter := returnValues[0].Interface()
	errmsg := ""
	if errInter != nil {
		errmsg = errInter.(error).Error()
	}
	// 然后发送出去
	server.sendResponse(sending, req, replyv.Interface(), codec, errmsg)
	server.freeRequest(req)
}
```

## 最后就是如果复用 net/http 包进行注册/监听

感觉这里有两个 方向一个Accept 就是链接来了怎么处理，直到返回非nil error

 一个就是ServeHttp应对单个的请求

```
// Accept accepts connections on the listener and serves requests
// for each incoming connection. Accept blocks until the listener
// returns a non-nil error. The caller typically invokes Accept in a
// go statement.
// accept 怎么接受链接 接受一次链接，就开一个线程
func (server *Server) Accept(lis net.Listener) {
	for {
		conn, err := lis.Accept()
		if err != nil {
			return
		}
		go server.ServeConn(conn)
	}
}
// ServeConn runs the server on a single connection.
// ServeConn 在单一的链接上
// ServeConn blocks, serving the connection until the client hangs up.
// ServeConn 会阻塞，服务这个链接直到客户端断开
// The caller typically invokes ServeConn in a go statement.
// caller 会调用serveconn 在一个go 程序里面
// ServeConn uses the gob wire format (see package gob) on the
// connection. To use an alternate codec, use ServeCodec.
// serveConn 使用 gob 写 在这个链接里面。 可以用别的
// See NewClient's comment for information about concurrent access.
func (server *Server) ServeConn(conn io.ReadWriteCloser) {
	buf := bufio.NewWriter(conn)
	srv := &gobServerCodec{
		rwc:    conn,
		dec:    gob.NewDecoder(conn),
		enc:    gob.NewEncoder(buf),
		encBuf: buf,
	}
	server.ServeCodec(srv)
}
// ServeCodec is like ServeConn but uses the specified codec to
// decode requests and encode responses.
// ServeCodec 和 ServeConn 类似但是使用 特定的codec 来编码和解码
func (server *Server) ServeCodec(codec ServerCodec) {
	sending := new(sync.Mutex)
	wg := new(sync.WaitGroup)
	// 会一直轮询直到链接被客户端断掉
	for {
		service, mtype, req, argv, replyv, keepReading, err := server.readRequest(codec)
		if err != nil {
			if debugLog && err != io.EOF {
				log.Println("rpc:", err)
			}
			if !keepReading {
				break
			}
			// send a response if we actually managed to read a header.
			if req != nil {
				server.sendResponse(sending, req, invalidRequest, codec, err.Error())
				server.freeRequest(req)
			}
			continue
		}
		wg.Add(1)
		// 这里call 会调用函数，然后发送出去
		go service.call(server, sending, wg, mtype, req, argv, replyv, codec)
	}
	// We've seen that there are no more requests.
	// Wait for responses to be sent before closing codec.
	wg.Wait()
	codec.Close()
}

// ServeHTTP implements an http.Handler that answers RPC requests.
// ServeHTTP 实现了一个http.Handler 用于回答RPC 请求
func (server *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// req Method 是CONNECT method
	if req.Method != "CONNECT" {
		// W.Header().Set() 设置http 头部
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		// 写入头部 之
		w.WriteHeader(http.StatusMethodNotAllowed)
		// 得先Connect 
		io.WriteString(w, "405 must CONNECT\n")
		return
	}
	// 这里Hijack 是直接取出tcp ，这里就是直接返回了，正常的链接会等待继续使用
	conn, _, err := w.(http.Hijacker).Hijack()
	if err != nil {
		log.Print("rpc hijacking ", req.RemoteAddr, ": ", err.Error())
		return
	} 
	// io.wRTIE 写入
	io.WriteString(conn, "HTTP/1.0 "+connected+"\n\n")
	// 存储Conn的
	server.ServeConn(conn)
}

```

## 后续的流程图

![Server](E:\Scattered-Study-Notes\go-rpc\Server.png)