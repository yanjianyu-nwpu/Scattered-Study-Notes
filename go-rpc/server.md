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