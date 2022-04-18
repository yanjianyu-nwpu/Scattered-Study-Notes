# net/rpc client

 

```
// Call represents an active RPC.
// Call 代表活跃的rpc
type Call struct {
	ServiceMethod string      // The name of the service and method to call.
	Args          interface{} // The argument to the function (*struct).
	Reply         interface{} // The reply from the function (*struct).
	Error         error       // After completion, the error status.
	Done          chan *Call  // Receives *Call when Go is complete.
}

// Client represents an RPC Client.
// client 待变一个rpc 客户端
// There may be multiple outstanding Calls associated
// with a single Client, and a Client may be used by
// multiple goroutines simultaneously.
// 可能有很多突出的调用 会被联系到单一的 客户端，多个协程调用。
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
// A ClientCodec implements writing of RPC requests and
// reading of RPC responses for the client side of an RPC session.
// 一个客户端ClientCodec 实现了读取rpc 请求和读取rpc 返回在rpc 会话的一端。
// The client calls WriteRequest to write a request to the connection
// and calls ReadResponseHeader and ReadResponseBody in pairs
// to read responses. 
// 客户端调用调用writeRequest 来写请求 并调用ReadResponseHeader 和 ReadResponseBody 成对读取返回值
//The client calls Close when finished with the
// connection. ReadResponseBody may be called with a nil
// argument to force the body of the response to be read and then
// discarded.
// 客户端会调用Close 当关闭，ReadResponseBody可能会用nil调用，强制读取response并丢弃
// See NewClient's comment for information about concurrent access.
type ClientCodec interface {
	WriteRequest(*Request, interface{}) error
	ReadResponseHeader(*Response) error
	ReadResponseBody(interface{}) error

	Close() error
}
// send
func (client *Client) send(call *Call) {
	client.reqMutex.Lock()
	defer client.reqMutex.Unlock()

	// Register this call.
	client.mutex.Lock()
	if client.shutdown || client.closing {
		client.mutex.Unlock()
		call.Error = ErrShutdown
		call.done()
		return
	}
	// 一个客户端会有一个seq，序列号好debug
	seq := client.seq
	client.seq++
	client.pending[seq] = call
	client.mutex.Unlock()

	// Encode and send the request.
	// 编码和发送请求
	client.request.Seq = seq
	client.request.ServiceMethod = call.ServiceMethod
	// 赋值情况
	// 编码请求
	err := client.codec.WriteRequest(&client.request, call.Args)
	if err != nil {
		client.mutex.Lock()
		call = client.pending[seq]
		delete(client.pending, seq)
		client.mutex.Unlock()
		if call != nil {
			call.Error = err
			call.done()
		}
	}
}

// input 输出
func (client *Client) input() {
	var err error
	var response Response
	for err == nil {
		response = Response{}
		// 读取返回值
		err = client.codec.ReadResponseHeader(&response)
		if err != nil {
			break
		}
		// 这里的拿到新的序列号
		seq := response.Seq
		client.mutex.Lock()
		call := client.pending[seq]
		delete(client.pending, seq)
		client.mutex.Unlock()

		switch {
		case call == nil:
			// We've got no pending call. 
			//That usually means that WriteRequest partially failed, and call was already removed; response is a server telling us about an error reading request body. We should still attempt to read error body, but there's no one to give it to.			
			// 常writeRequset 可能获取失败 仍旧要尝试读取body
			err = client.codec.ReadResponseBody(nil)
			if err != nil {
				err = errors.New("reading error body: " + err.Error())
			}
		case response.Error != "":
			// We've got an error response. Give this to the request;
			// any subsequent requests will get the ReadResponseBody
			// error if there is one.
			// 如果接收到error 返回值，
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

