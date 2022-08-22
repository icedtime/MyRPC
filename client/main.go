package main

import (
	"fmt"
	"myrpctest/client/impl"
)

func main() {
	//从服务注册中心查询要调用的服务的ip和端口
	//现在的话没有服务注册中心，所以先写死，直接给定服务的IP和端口
	//现在默认tcp
	client := impl.NewClient("127.0.0.1", 7777)
	//连接服务器获得一个conn，然后封装conn为connection
	//TODO 怎么实现直接连接服务，连接服务之后CALL输入方法名后直接调用
	//这里暂时实现，在生成connection的时候存储服务名和方法名
	connect, err := client.NewConnect() // 对特定服务 建立连接
	if err != nil {
		fmt.Println(err)
		return
	}
	defer connect.Stop()

	req := impl.SendMessage{
		Data: "爱已随风起",
	}
	resp := impl.NewResponse()
	err = connect.Call("Test", "SayHello", &req, resp) // 调用某个服务
	if err != nil {
		fmt.Println()
	}

	fmt.Println(resp.Msg.Data)
	fmt.Println("开始关闭conn！")
}
