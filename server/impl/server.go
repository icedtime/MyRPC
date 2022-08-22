package impl

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"reflect"
	"server/message"
	"time"
)

type Server struct {
	Name      string
	IPVersion string
	IP        string
	Port      int

	isClose chan struct{}
	//服务器配置
	ServerOption *Options
	//注册服务
	ServiceMap map[string]*Service
	//函数调用

	//允许的最大连接数
	MaxConnCount int
	//当前的连接数
	ConnCount int
	//当前服务器的链接
	ConnList map[string]net.TCPConn
}

//规定提供的可以RPC调用的函数必须是以下形式
//TODO 加入ctx
type HandleFunc func(req *message.Request, res *message.Response) error

func NewServer(name string) *Server {
	return &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      7777,

		ConnCount: 0,
		//初始化一个服务器
		ServerOption: DefaultOptions(),
		//本地服务注册
		ServiceMap: map[string]*Service{},
		ConnList:   map[string]net.TCPConn{},
	}
}

//包括三部分：
//读取配置options
//服务注册，包括本地和服务注册中心
//服务跑起来
func (s *Server) Serve() error {
	//一、读取配置
	//TODO 我这里不用读取配置，因为默认位置，我是直接给定，然后赋值给服务器
	//之所以有这一部分是因为，一旦配置项很多，在初始化的时候，就要传入很多很多值，所以需要写一个新的初始化函数
	//包括获取listenner，并封装listenner
	addr, _ := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", s.IP, s.Port))
	ls, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}
	s.ServerOption.ls = ls

	//二、服务注册

	//三、服务跑起来，开启子协程监听端口，不断接受请求conn，并处理
	return s.Start()

	//start里面开启了额外的协程，如果不用这句阻塞，直接整个结束
	select {}
}

func (s *Server) Start() error {
	fmt.Println("server start ...")
	// fmt.Println("打印一下注册了的服务!")
	// for k, _ := range s.ServiceMap {
	// 	fmt.Println(k)
	// }

	// go func() {
	// 	//获取服务器IP:port
	// 	fmt.Println("resolve TCP address...")
	// 	TcpAddr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
	// 	if err != nil {
	// 		fmt.Println("resolve tcp addr error :", err)
	// 		return
	// 	}
	// 	fmt.Println("resolve TCP address success!")
	// 	fmt.Println("accept listenner...")
	// 	//network一般指TCP、HTTP的版本
	// 	//监听端口
	//listenner, err := net.ListenTCP(s.IPVersion, TcpAddr)
	// 	if err != nil {
	// 		fmt.Println("listen tcp addr error :", err)
	// 		return
	// 	}
	// 	fmt.Println("listen TCP addr success!")

	// 	//不停地获取链接conn
	// 	for {
	// 		fmt.Println("accept conn start...")
	// 		//TODO conn不能只做一次调用，要多次,加循环就可以了
	// 		//TODO conn怎么判断已经断开？
	// 		conn, err := listenner.AcceptTCP()
	// 		if err != nil {
	// 			fmt.Println("accept conn error :", err)
	// 			return
	// 		}
	// 		fmt.Println("accept conn success!")

	// 		//开启另外的协程处理请求
	// 		go s.DealConn(conn)
	// 	}
	// }()
loop:
	for {
		select {
		case <-s.ServerOption.ctx.Done(): // 检查是否需要退出服务
			break loop
		default:
			//这里其实还应该适应其他的协议,太麻烦，只写TCP
			conn, err := s.ServerOption.ls.AcceptTCP() // 获取一个链接
			if err != nil {
				log.Println(err)
				continue
			}

			//TODO 不知道干嘛用的
			// if s.options.Trace {
			// 	log.Println("connect: ", accept.RemoteAddr())
			// }

			go s.DealConn(conn) // 开一个协程去处理 该 链接
		}

	}
	return nil
}

func (s *Server) Stop() {

}

func (s *Server) Call(req *message.Request, res *message.Response) error {
	serverName := req.ServerName
	methodName := req.MethodName

	service, ok := s.ServiceMap[serverName]
	if !ok {
		return errors.New("没有找到服务！")
	}

	methodType, ok := service.MethodType[methodName]
	if !ok {
		return errors.New("没有找到方法！")
	}
	fn := methodType.Method.Func
	returnValue := fn.Call([]reflect.Value{service.RefVal, reflect.ValueOf(req), reflect.ValueOf(res)})
	errInterface := returnValue[0].Interface()
	if errInterface != nil {
		return errInterface.(error)
	}
	fmt.Println("111111111111", res.Msg.Data)
	return nil
}

//本地注册服务
func (s *Server) RegisterUseName(service interface{}, name string) error {
	return s.RegisterLocalService(service, true, name)
}

func (s *Server) RegisterWithoutName(Service interface{}) error {
	return s.RegisterLocalService(Service, false, "")
}

func (s *Server) RegisterLocalService(service interface{}, useName bool, name string) error {
	ser, err := NewService(service, useName, name)
	if err != nil {
		return err
	}
	s.ServiceMap[ser.Name] = ser
	return nil
}

//服务注册中心注册
//启动服务器之前就把所有的服务注册，启动之后不注册了
func (s *Server) RegisterRemoteService(serverName []string) {
	//获取discovery中心的地址
	disAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", s.ServerOption.DiscoveryServerIp, s.ServerOption.DiscoveryServerPort))
	if err != nil {
		fmt.Println("获取discovery地址出错！")
		return
	}
	fmt.Println("获取discovery地址成功！")

	//与dicovery建立连接
	conn, err := net.DialTCP("tcp", nil, disAddr)
	if err != nil {
		fmt.Println("连接discovery失败！")
		return
	}
	fmt.Println("连接discovery成功！")

	//TODO应该封装为一个方法
	//因为目前我的服务注册中心，注册其实就是一个serverName后面跟提供服务的服务器IP
	addServerMsg := message.NewAddServerMsg(serverName)
	DisRequest := message.NewDisRequest("addServer", addServerMsg)
	//序列化
	jsonDisReq, err := json.Marshal(DisRequest)
	if err != nil {
		fmt.Println("序列化失败！")
		return
	}
	fmt.Println("序列化成功!")
	//发送给服务注册中心
	conn.Write(jsonDisReq)

	//TODO 这里应该接收服务注册中心的返回的注册成功与否的信息
	//如果错误就要抛出错误
	//现在不写

	//TODO 感觉conn应该封装一下
	//封装完之后，开启心跳
	//我不知道这样的逻辑行不行
	//注册完之后就开启心跳
	go s.HeartBeat()
	//主协程不能挂，挂了子协程也会挂
	select {}
}

//TODO 做到读写分离
func (s *Server) DealConn(conn *net.TCPConn) {
	//网络不可靠的时候的处理
	// defer func() {
	// 	// 网络不可靠
	// 	if err := recover(); err != nil {
	// 		utils.PrintStack()
	// 		log.Println("Recover Err: ", err)
	// 	}
	// }()

	//服务器连接数增加
	s.ConnCount++
	addr := conn.RemoteAddr().String()
	s.ConnList[addr] = conn
	//用来处理conn的关闭，还有server连接情况
	defer func() {
		//当链接断开之后，修改
		s.ConnCount--
		delete(s.ConnList, addr)
		err := conn.Close()
		if err != nil {
			fmt.Println("conn 关闭出错:", err)
		}
	}()

	//conn初始化，将conn转为TCPConn
	// if tcpConn, ok := conn.(*net.TCPConn); ok {
	// 	//是对tcp的一些初始化
	// 	//TODO 了解一下TCP的初始化事项,是否是要保持长连接，linger等
	// }

	//TODO 实现读写分离
	//不断从conn读出数据
	for {
		//读出序列化的数据
		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Print("read msg from conn error :", err)
			return
		}
		//反序列化得到request
		//req的Msg没有初始化，却不报错？
		var req message.Request
		err = json.Unmarshal(buf[:cnt], &req)
		if err != nil {
			fmt.Println("反序列化失败:", err)
			return
		}
		fmt.Println("反序列化成功！")
		fmt.Println("接收到的服务名和方法名为:", req.ServerName, req.MethodName)
		//测试一下是否真的发送数据过来，并接收成功
		//接受成功了
		//fmt.Println(req.Msg.Data)

		//response
		//TODO 使用服务名和方法名调用方法
		//暂时写死
		//报空指针异常，很简单，上面res还嵌套着一层msg，那个结构体没有初始化
		res := message.NewResponse()
		// test := NewTest()
		// err = test.SayHello(&req, res)
		//函数类型就不要传指针了
		err = s.Call(&req, res)
		if err != nil {
			fmt.Println("调用函数出错:", err)
			return
		}
		fmt.Println(res.Msg.Data)

		//序列化response
		jsonRes, err := json.Marshal(res)
		if err != nil {
			fmt.Println("序列化失败！")
			return
		}
		fmt.Println("序列化成功！")
		_, err = conn.Write([]byte(jsonRes))
		if err != nil {
			fmt.Println("write back buf err", err)
			return
		}

	}

	//当链接断开之后，修改
	s.ConnCount--
	delete(s.ConnList, addr)
}

func (s *Server) HeartBeat() {
	for {
		select {
		case <-s.isClose:
			return
		case <-time.After(time.Second * s.ServerOption.readTimeout):
			//从服务注册中心那里获得连接池
		}
	}
}
