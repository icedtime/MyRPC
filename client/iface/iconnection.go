package iface

import "net"

type Iconnection interface {
	//连接远程服务器
	Call(serverName string, request interface{}, response interface{}) error
	//停止链接，结束当前链接的工作
	Stop()
	//获取当前链接的绑定socket conn
	//socket套接字
	GetTCPConnection() *net.Conn
	//获取当前链接的链接ID
	GetConnID() uint32
	//获取远程客户端的TCP状态：IP:port
	RemoteAddr() net.Addr
	//发送数据，将数据发送到远方客户端
	Send(data []byte) error
}
