package transport

import (
	"errors"
	"fmt"
	"net"
)

type transport struct {
	trMap map[Protocol]genTransport
}

type genTransport func(addr string) (net.Listener, error)

//初始化
var Transport = &transport{
	trMap: map[Protocol]genTransport{},
}

func (t *transport) registry(trType Protocol, genTr genTransport) {
	t.trMap[trType] = genTr
}

//最重要的功能，生成对应传输协议的获取listenner的函数
func (t *transport) Gen(trType Protocol, addr string) (net.Listener, error) {
	//通过配置中的传输协议获得方法
	gen, ok := t.trMap[trType]
	if !ok {
		return nil, errors.New(fmt.Sprintf("%s协议未找到！", trType))
	}

	return gen(addr)
}
