package impl

type Request struct {
	ServerName string
	MethodName string
	//有时候在想这里这个是指针吧，把指针发过去真的可以吗？
	//为什么不是发数据？
	Msg interface{}
}

func NewRequest(serverName, methodName string, msg interface{}) *Request {
	return &Request{
		ServerName: serverName,
		MethodName: methodName,
		Msg:        msg,
	}
}
