package message

type Request struct {
	ServerName string
	MethodName string
	//有时候在想这里这个是指针吧，把指针发过去真的可以吗？
	//为什么不是发数据？
	Msg *SendMessage
}

func NewRequest(serverName, methodName string, msg *SendMessage) *Request {
	return &Request{
		ServerName: serverName,
		MethodName: methodName,
		Msg:        msg,
	}
}

//该结构用于服务端和服务注册中心、客户端和服务注册中心的交互，服务端和客户端的交互用其他的形式
type DisRequest struct {
	//对于服务端提供两种:heartBeat和addServer
	//对于客户端提供一种：searchServer
	ServerName string
	Msg        interface{}
}

func NewDisRequest(serverName string, msg interface{}) *DisRequest {
	return &DisRequest{
		ServerName: serverName,
		Msg:        msg,
	}
}

type HeartBeatMsg struct {
	//用于发送当前服务的conn链接数量，实现负载均衡
	ConnCount int
}

func NewHeartBeatMsg(connCount int) *HeartBeatMsg {
	return &HeartBeatMsg{
		ConnCount: connCount,
	}
}

type AddServerMsg struct {
	//注册一组服务
	ServerNames []string
}

func NewAddServerMsg(serverNames []string) *AddServerMsg {
	return &AddServerMsg{
		ServerNames: serverNames,
	}
}

type SearchServerMsg struct {
	//查找一个服务
	ServerName string
}

func NewSearchServerMsg(serverName string) *SearchServerMsg {
	return &SearchServerMsg{
		ServerName: serverName,
	}
}
