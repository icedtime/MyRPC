package transport

//在网络传输中并不是只会用tcp，还可以用其他的
//这个protocol就是为了针对不同传输协议写的
//目前只写了tcp
//TODO 写其他协议
//有什么用？
//socket编程的步骤，获取addr，监听addr获得listenner，从listenner中获取conn
//不同协议不同的那一步就是获取listenner，所以要根据不同协议生成不同的获取listenner方法
type Protocol string

const (
	TCP Protocol = "tcp"
)

func (p Protocol) String() string {
	return string(p)
}
