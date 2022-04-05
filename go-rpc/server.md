 package main

import (
    "fmt"
    "io"
    "net"
    "net/http"
    "net/rpc"
)

type Watcher int

func (w *Watcher) GetInfo(arg int, result *int) error {
    *result = 1
    return nil
}

func main() {

    http.HandleFunc("/ghj1976", Ghj1976Test)
    
    watcher := new(Watcher)
    rpc.Register(watcher)
    rpc.HandleHTTP()
    
    l, err := net.Listen("tcp", ":1234")
    if err != nil {
        fmt.Println("监听失败，端口可能已经被占用")
    }
    fmt.Println("正在监听1234端口")
    http.Serve(l, nil)
}

func Ghj1976Test(w http.ResponseWriter, r *http.Request) {
    io.WriteString(w, "<html><body>ghj1976-123</body></html>")
}

 

客户端代码：

package main

import (
    "fmt"
    "net/rpc"
)

func main() {
    client, err := rpc.DialHTTP("tcp", "127.0.0.1:1234")
    if err != nil {
        fmt.Println("链接rpc服务器失败:", err)
    }
    var reply int
    err = client.Call("Watcher.GetInfo", 1, &reply)
    if err != nil {
        fmt.Println("调用远程服务失败", err)
    }
    fmt.Println("远程服务返回结果：", reply)
}



# 标记



suitableMethods

```
func suitableMethods(typ reflect.Type, reportErr bool) map[string]*methodType {
	//感觉是用于获取函数
	methods := make(map[string]*methodType)
	// 所有的函数都被添加
	for m := 0; m < typ.NumMethod(); m++ {
	
		//  找到对象
		method := typ.Method(m)
		mtype := method.Type
		mname := method.Name
		// Method must be exported.
		
		if method.PkgPath != "" {
			continue
		}
		// Method needs three ins: receiver, *args, *reply.
		// rpc 指定是有三个参数，第一个就是结构体，要是（l *server）的结构体
		if mtype.NumIn() != 3 {
			if reportErr {
				log.Printf("rpc.Register: method %q has %d input parameters; needs exactly three\n", mname, mtype.NumIn())
			}
			continue
		}
		// 第二个是要是 指针,就是结构体要附在指针上
		// First arg need not be a pointer.
		argType := mtype.In(1)
		if !isExportedOrBuiltinType(argType) {
			if reportErr {
				log.Printf("rpc.Register: argument type of method %q is not exported: %q\n", mname, argType)
			}
			continue
		}
		// 第二个入参也必须是 指针
		// Second arg must be a pointer.
		replyType := mtype.In(2)
		if replyType.Kind() != reflect.Ptr {
			if reportErr {
				log.Printf("rpc.Register: reply type of method %q is not a pointer: %q\n", mname, replyType)
			}
			continue
		}
		// 返回值也必须是可导出的
		// Reply type must be exported.
		if !isExportedOrBuiltinType(replyType) {
			if reportErr {
				log.Printf("rpc.Register: reply type of method %q is not exported: %q\n", mname, replyType)
			}
			continue
		}
		// Method needs one out.
		// method 需要一个输出
		if mtype.NumOut() != 1 {
			if reportErr {
				log.Printf("rpc.Register: method %q has %d output parameters; needs exactly one\n", mname, mtype.NumOut())
			}
			continue
		}
		// The return type of the method must be error.
		// 输出必须是error, k,
		if returnType := mtype.Out(0); returnType != typeOfError {
			if reportErr {
				log.Printf("rpc.Register: return type of method %q is %q, must be error\n", mname, returnType)
			}
			continue
		}
		methods[mname] = &methodType{method: method, ArgType: argType, ReplyType: replyType}
	}
	return methods
}
```



```
// 用于用注册 注册函数到 
func (server *Server) register(rcvr interface{}, name string, useName bool) error {
	s := new(service)
	s.typ = reflect.TypeOf(rcvr)
	s.rcvr = reflect.ValueOf(rcvr)
	// reflect.Indirect返回指向对象的类型
	sname := reflect.Indirect(s.rcvr).Type().Name()

	log.Println(sname)
	if useName {
		sname = name
	}
	if sname == "" {
		s := "rpc.Register: no service name for type " + s.typ.String()
		log.Print(s)
		return errors.New(s)
	}
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
	// 导入新的 sync 包
	if _, dup := server.serviceMap.LoadOrStore(sname, s); dup {
		return errors.New("rpc: service already defined: " + sname)
	}
	return nil
}
```

这里主要获取和校验所有的函数

符合所有的要求

然后再注册端口  

然后将端口注册 HandleHTTP()



会注册两个端口

DefaultRPCPath, DefaultDebugPath

DefaultRPCPath  = "/_goRPC_"

  DefaultDebugPath = "/debug/rpc"

然后这里/goRPC 这里路径是用于method 的使用

/debug/src 是可以访问，这里是会讲述怎么用rpc method





src/net/rpc/server.go

这里是 http.Handle 



然后再 go/SRC/NEt/hpp 包这里 server.go 文件

然后这里 有个默认实例 Hanler

```
func Handle(pattern string, handler Handler) { DefaultServeMux.Handle(pattern, handler) }
```



这里 Handler 是个 interface

```
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}
```

这里是实现一个 ServeHttp  的interface 然后这个Handle 的定义

```
// Handle registers the handler for the given pattern.
// Handle注册 指定的pattern'
// If a handler already exists for pattern, Handle panics.、
// 如果重复会panic
func (mux *ServeMux) Handle(pattern string, handler Handler) {
// 上锁
	mux.mu.Lock()
	defer mux.mu.Unlock()

	if pattern == "" {
		panic("http: invalid pattern")
	}
	if handler == nil {
		panic("http: nil handler")
	}
	if _, exist := mux.m[pattern]; exist {
		panic("http: multiple registrations for " + pattern)
	}

	if mux.m == nil {
		mux.m = make(map[string]muxEntry)
	}
	e := muxEntry{h: handler, pattern: pattern}
	mux.m[pattern] = e
	if pattern[len(pattern)-1] == '/' {
		mux.es = appendSorted(mux.es, e)
	}

	if pattern[0] != '/' {
		mux.hosts = true
	}
}
```

然后rpc 这里实现的 这

```
// ServeHTTP implements an http.Handler that answers RPC requests.
// ServeHttp 实现了 http。Hanler 去返回rpc 请求
func (server *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != "CONNECT" {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusMethodNotAllowed)
		io.WriteString(w, "405 must CONNECT\n")
		return
	}
	conn, _, err := w.(http.Hijacker).Hijack()
	if err != nil {
		log.Print("rpc hijacking ", req.RemoteAddr, ": ", err.Error())
		return
	}
	io.WriteString(conn, "HTTP/1.0 "+connected+"\n\n")
	server.ServeConn(conn)
}
```



ＳＥＮＤＲｅｓｐｏｎｃｅ发送

```
func (server *Server) sendResponse(sending *sync.Mutex, req *Request, reply interface{}, codec ServerCodec, errmsg string) {
／／　先得到ｒｅｐｓｐｏｎｓｅ
	resp := server.getResponse()
	// Encode the response header
	／／　得到ｓｅｒｖｉｄｅ　ｍｔｈｏｓ这样，
	resp.ServiceMethod = req.ServiceMethod
	if errmsg != "" {
		resp.Error = errmsg
		reply = invalidRequest
	}
	／／　复制序列号
	resp.Seq = req.Seq
	sending.Lock()
	// 进行压缩
	err := codec.WriteResponse(resp, reply)
	if debugLog && err != nil {
		log.Println("rpc: writing response:", err)
	}
	sending.Unlock()
	server.freeResponse(resp)
}
```

// 在ｍ　＊ｍｔｈｏｄＴｙｐｅ　中会维护一个ｎｍｂＣａｌｌｓ数目

这里世界的ｃａｌｌ　会指定一个　ｃａｌｌ

```
func (s *service) call(server *Server, sending *sync.Mutex, wg *sync.WaitGroup, mtype *methodType, req *Request, argv, replyv reflect.Value, codec ServerCodec) {
	if wg != nil {
		defer wg.Done()
	}
	// 上锁记录
	mtype.Lock()
	mtype.numCalls++
	mtype.Unlock()
	// 找到ｆｕｎ
	function := mtype.method.Func
	
	// Invoke the method, providing a new value for the reply.
	// 调用方法，提供新的值
	这里这么调用，
	returnValues := function.Call([]reflect.Value{s.rcvr, argv, replyv})
	// The return value for the method is an error.
	找到是否结果是ｅｒｒｏｒ　
	errInter := returnValues[0].Interface()
	errmsg := ""
	if errInter != nil {
		errmsg = errInter.(error).Error()
	}
	server.sendResponse(sending, req, replyv.Interface(), codec, errmsg)
	server.freeRequest(req)
}
```

ｓｅｒｖｅＣｏｎｎDE 

在单链接的情况下。会保存链接直到客户端挂起

# Client文件

```
// Client represents an RPC Client.
// 客户端必须待变一个RPC客户端
// There may be multiple outstanding Calls associated
// with a single Client, and a Client may be used by
// multiple goroutines simultaneously.
// 
type Client struct {
	codec ClientCodec

	reqMutex sync.Mutex // protects following
	request  Request

	mutex    sync.Mutex // protects following
	seq      uint64
	pending  map[uint64]*Call
	closing  bool // user has called Close
	shutdown bool // server has told us to stop
}
```

然后这里 感觉就是客户段，有个序列号和锁搞没有竞态环节



```


// A ClientCodec implements writing of RPC requests and
// reading of RPC responses for the client side of an RPC session.
// The client calls WriteRequest to write a request to the connection
// and calls ReadResponseHeader and ReadResponseBody in pairs
// to read responses. The client calls Close when finished with the
// connection. ReadResponseBody may be called with a nil
// argument to force the body of the response to be read and then
// discarded.
// ClientCodec 实现了 写rpc 请求并 读取rpc 
// 客户端调用writeRequest 用于 写request
// 调用 readResponseHeader 和 readResponseBody 这一对函数读response
// 客户端关闭 当connection 关闭的时候，ReadResponseBody 可能是nil 来强制读取response
// 并丢弃，从NewClient 的注释来 获取关于并发控制
// See NewClient's comment for information about concurrent access.
type ClientCodec interface {
	WriteRequest(*Request, interface{}) error
	ReadResponseHeader(*Response) error
	ReadResponseBody(interface{}) error

	Close() error
}
```



然后send 函数值 

```
func (client *Client) send(call *Call) {
	// 上锁这个req的锁
	client.reqMutex.Lock()
	defer client.reqMutex.Unlock()

	// Register this call.
	// 这个req 感觉是统计维护 序列号
	client.mutex.Lock()
	if client.shutdown || client.closing {
		client.mutex.Unlock()
		call.Error = ErrShutdown
		call.done()
		return
	}
	seq := client.seq
	client.seq++
	client.pending[seq] = call
	client.mutex.Unlock()

	// Encode and send the request.
	// 编码并且发送请求请求，
	client.request.Seq = seq
	// 记录一下
	client.request.ServiceMethod = call.ServiceMethod
	log.Print(client.request)
	//  压缩请求 
	err := client.codec.WriteRequest(&client.request, call.Args)
	// 这里就是
	if err != nil {
		client.mutex.Lock()
		call = client.pending[seq]
		delete(client.pending, seq)
		client.mutex.Unlock()
		if call != nil {
			call.Error = err
			call.done()·
		}
	}
}
```

input 函数

```
// 
func (client *Client) input() {
	var err error
	var response Response
	for err == nil {
		response = Response{}
		// 解压 response
		err = client.codec.ReadResponseHeader(&response)
		if err != nil {
			break
		}
		seq := response.Seq
		client.mutex.Lock()
		call := client.pending[seq]
		delete(client.pending, seq)
		client.mutex.Unlock()

		switch {
		case call == nil:
			// We've got no pending call. That usually means that
			// WriteRequest partially failed, and call was already
			// removed; response is a server telling us about an
			// error reading request body. We should still attempt
			// to read error body, but there's no one to give it to.
			err = client.codec.ReadResponseBody(nil)
			if err != nil {
				err = errors.New("reading error body: " + err.Error())
			}
		case response.Error != "":
			// We've got an error response. Give this to the request;
			// any subsequent requests will get the ReadResponseBody
			// error if there is one.
			call.Error = ServerError(response.Error)
			err = client.codec.ReadResponseBody(nil)
			if err != nil {
				err = errors.New("reading error body: " + err.Error())
			}
			call.done()
		default:
			err = client.codec.ReadResponseBody(call.Reply)
			if err != nil {
				call.Error = errors.New("reading body " + err.Error())
			}
			call.done()
		}
	}
	// Terminate pending calls.
	client.reqMutex.Lock()
	client.mutex.Lock()
	client.shutdown = true
	closing := client.closing
	if err == io.EOF {
		if closing {
			err = ErrShutdown
		} else {
			err = io.ErrUnexpectedEOF
		}
	}
	for _, call := range client.pending {
		call.Error = err
		call.done()
	}
	client.mutex.Unlock()
	client.reqMutex.Unlock()
	if debugLog && err != io.EOF && !closing {
		log.Println("rpc: client protocol error:", err)
	}
}
```

```
// NewClient returns a new Client to handle requests to the
// set of services at the other end of the connection.
// It adds a buffer to the write side of the connection so
// the header and payload are sent as a unit.
// NewClient 返回一个新的Client 会添加bufuer 哟关于写单侧写 所以头部和负载会作为一个单位发送
// The read and write halves of the connection are serialized independently,
//  读写 在链接两端 是独立的 而然 每个半边是可以接受并发的 所以要并发安全
// so no interlocking is required. However each half may be accessed
// concurrently so the implementation of conn should protect against
// concurrent reads or concurrent writes.
func NewClient(conn io.ReadWriteCloser) *Client {
	encBuf := bufio.NewWriter(conn)
	client := &gobClientCodec{conn, gob.NewDecoder(conn), gob.NewEncoder(encBuf), encBuf}
	return NewClientWithCodec(client)
}
```

```
type gobClientCodec struct {
	rwc    io.ReadWriteCloser
	dec    *gob.Decoder
	enc    *gob.Encoder
	encBuf *bufio.Writer
}
// 用于可客户端压缩
```

