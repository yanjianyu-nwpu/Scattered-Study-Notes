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

然后rpc 这里实现的

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



