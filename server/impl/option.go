package impl

import (
	"context"
	"net"
	"server/transport"
	"time"
)

type Options struct {
	DiscoveryServerIp   string
	DiscoveryServerPort int

	proto transport.Protocol
	//封装一个listenner
	ls *net.TCPListener
	//服务器最大连接数，用于熔断
	MaxConnCount int
	//超时时间，用于心跳
	readTimeout time.Duration

	//控制服务退出的context
	ctx context.Context
	//TODO,如果之后要写https
	//还可以写公钥和私钥
}

//固定的，不能由其他方法修改
//固定默认配置
func DefaultOptions() *Options {
	return &Options{
		DiscoveryServerIp:   "127.0.0.1",
		DiscoveryServerPort: 8888,
		//默认TCP连接
		proto:        transport.TCP,
		MaxConnCount: 512,
		readTimeout:  time.Minute * 3,
		ctx:          context.Background(),
	}
}
