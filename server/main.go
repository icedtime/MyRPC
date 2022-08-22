package main

import (
	"fmt"
	"server/impl"
)

func main() {
	//初始化服务器
	server := impl.NewServer("服务器一号！")
	//本地注册方法
	err := server.RegisterUseName(&impl.Test{}, "Test")
	if err != nil {
		fmt.Println("本地注册服务失败！")
	}
	fmt.Println("本地注册服务成功！")
	//服务注册中心注册方法
	err = impl.RegisterRemoteService([]string{"Test"})
	if err != nil {

	}
	fmt.Println()
	//启动服务器
	server.Serve()
	//TODO 写个优雅的关闭
}
