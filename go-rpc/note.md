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

```
// Register publishes in the server the set of methods of the
// receiver value that satisfy the following conditions:
// Register 发布
//	- exported method of exported type
//	- two arguments, both of exported type
//	- the second argument is a pointer
//	- one return value, of type error
// It returns an error if the receiver is not an exported type or has
// no suitable methods. It also logs the error using package log.
// The client accesses each method using a string of the form "Type.Method",
// where Type is the receiver's concrete type.
func (server *Server) Register(rcvr interface{}) error {
	return server.register(rcvr, "", false)
}
```

