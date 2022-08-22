package impl

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

type Connection struct {
	Conn   *net.TCPConn
	ConnID uint32
	//当前链接的状态
	isClosed bool

	//用于链接处理业务的函数API
	//handleAPI ziface.HandleFunc

	//通知当前链接已经退出/停止用的channel
	ExitChan chan bool

	//暂时这样处理，将服务名封装到connection
	//坏处就是一个连接只能进行 一次远程调用
	//TODO 让一个conn可以调用多个服务
	//ServerName string
	//MethodName string
}

func NewConnection(conn *net.TCPConn, connID uint32) *Connection {
	return &Connection{
		Conn:     conn,
		ConnID:   connID,
		isClosed: false,
		ExitChan: make(chan bool, 1),
	}
}

func (c *Connection) Call(serverName, methodName string, sendMessage *SendMessage, response *Response) error {
	//将请求信息和服务名封装之后发过去
	request := NewRequest(serverName, methodName, sendMessage)
	//序列化
	jsonReq, err := json.Marshal(request)
	if err != nil {
		return err
	}
	fmt.Println("序列化成功！")

	//发送数据
	_, err = c.Conn.Write(jsonReq)
	if err != nil {
		return err
	}
	fmt.Println("发送成功！")

	//接受数据返回值
	buf := make([]byte, 512)
	cnt, err := c.Conn.Read(buf)
	if err != nil {
		return err
	}
	fmt.Println("接受从客户端来的数据成功：", err)
	//反序列化
	//反序列化要加上长度 ，不然会报错
	err = json.Unmarshal(buf[:cnt], &response)
	if err != nil {
		return err
	}
	fmt.Println("反序列化成功！")

	fmt.Printf("server call back : %s", response.Msg.Data)
	time.Sleep(1 * time.Second)
	return nil
}

//停止链接，结束当前链接的工作
func (c *Connection) Stop() {
	fmt.Println("Conn Stop.. ConnID=", c.ConnID)

	//先判断链接是否关闭
	if c.isClosed {
		return
	}

	c.isClosed = true
	c.Conn.Close()
	//要记住回收资源
	close(c.ExitChan)
}

//获取当前链接的绑定socket conn
//socket套接字
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

//获取当前链接的链接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

//获取远程客户端的TCP状态：IP:port
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

//发送数据，将数据发送到远方客户端
// func (c *Connection) Send(data []byte) error {

// }
