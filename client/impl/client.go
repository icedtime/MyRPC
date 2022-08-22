package impl

import (
	"errors"
	"fmt"
	"net"
	"time"
)

type Client struct {
	//要连接的服务器IP和端口
	IP   string
	Port int
}

func NewClient(IP string, port int) *Client {
	return &Client{
		IP:   IP,
		Port: port,
	}
}

//TODO 返回TCPConn
func (c *Client) NewConnect() (*Connection, error) {
	time.Sleep(3 * time.Second)
	fmt.Println("客户端开始获取conn")
	//TODO 返回TCPConn
	//TCP有tcp、tcp4和tcp6
	//laddr *net.TCPAddr:本地地址，通常为nil，
	//raddr *net.TCPAddr：服务器地址
	fmt.Println("开始拼接raddr...")
	raddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", c.IP, c.Port))
	if err != nil {
		//%v按值原来形式输出
		newError := errors.New(fmt.Sprintf("拼接raddr出错:%v", err))
		return nil, newError
	}
	fmt.Println("拼接raddr成功！")

	fmt.Println("开始获取conn...")
	conn, err := net.DialTCP("tcp", nil, raddr)
	if err != nil {
		newError := errors.New(fmt.Sprintf("获取conn出错:%v", err))
		return nil, newError
	}
	fmt.Println("客户端获取conn成功！")

	//TODO 设置一个工具类，用于随机生成不重复的connID
	//现在写死
	connection := NewConnection(conn, 113)
	return connection, nil
}
