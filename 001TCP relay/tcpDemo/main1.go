package main

import (
	"bytes"
	"bufio"
	"fmt"
	"net"
	"time"
)

var ListernTCP = "0.0.0.0:8001"
var DailTCP = "43.139.232.236:8003"

// TCP server端

// 处理函数
func processServer(conn net.Conn, c1 chan []byte, c2 chan []byte) {
	defer conn.Close() // 关闭连接
	go func() {
		for {
			reader := bufio.NewReader(conn)
			var buf [40960]byte
			n, err := reader.Read(buf[:]) // 使用带缓冲方式读取数据
			if err != nil {
				fmt.Println("read from client failed, err:", err)
				break
			}
			//recvStr := string(buf[:n])
			fmt.Println("6000processServer：", string(buf[:n]))
			c1 <- buf[:n]
		}

	}()
	//client:
	//	for {
	//		select {
	//		case v := <-c2:
	//			conn.Write(v) // 发送数据
	//		default:
	//			goto client
	//		}
	//	}
	go func() {
		for {
			select {
			case v := <-c2:
				// 消息体长度和协议中长度不一致，导致一直等待
				req := bytes.NewBuffer(nil)
				req.Write(v)
				req.Write([]byte{'1','2','3','4','5','6','7','8','9','9','9','"','}'})
				conn.Write(req.Bytes()) // 发送数据
				fmt.Println("write server")
				fmt.Println(":",string(req.Bytes()))
				// conn.Close()
			default:
				continue
			}
		}
	}()
	for {

	}
}
func processClient(conn net.Conn, c1 chan []byte, c2 chan []byte) {
	defer conn.Close() // 关闭连接
	go func() {
		for {
			reader := bufio.NewReader(conn)
			var buf [12800]byte
			n, err := reader.Read(buf[:]) // 使用带缓冲方式读取数据
			if err != nil {
				fmt.Println("read from client failed, err:", err)
				break
			}
			fmt.Println("processClient：", string(buf[:n]))
			// 长度不够导致浏览器一直等待
			
			c2 <- buf[:n]
		}
	}()
	//client2:
	//	for {
	//		select {
	//		case v := <-c1:
	//			conn.Write(v)
	//		default:
	//			goto client2
	//		}
	//	}
	go func() {
		for {
			select {
			case v := <-c1:
				conn.Write(v)
				fmt.Println("write client")
			default:
				continue
			}
		}
	}()
	for {

	}
}
func main() {
	listen, err := net.Listen("tcp", ListernTCP)
	if err != nil {
		fmt.Println("listen failed, err:", err)
		return
	}
	time.Sleep(1 * time.Second)
	conn2, err := net.Dial("tcp", DailTCP)
	if err != nil {
		fmt.Println("err :", err)
		return
	}
	defer conn2.Close() // 关闭连接

	ch1 := make(chan []byte, 40) // 主机A服务端写入，客户端读取
	ch2 := make(chan []byte, 40) // 客户端写入收到的信息，服务器端读取

	for {
		conn, err := listen.Accept() // 建立连接
		if err != nil {
			fmt.Println("accept failed, err:", err)
			continue
		}
		go processServer(conn, ch1, ch2) // 启动一个goroutine处理连接
		go processClient(conn2, ch1, ch2)
	}
}
