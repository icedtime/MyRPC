package impl

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"

	"github.com/go-redis/redis"
)

type Discovery struct {
	RedisClient *redis.Client
	IPVersion   string
	IP          string
	Port        int

	//加入一个负载均衡算法

}

func NewDiscovery(redisClient *redis.Client) *Discovery {
	return &Discovery{
		RedisClient: redisClient,
		IPVersion:   "tcp4",
		IP:          "0.0.0.0",
		Port:        8888,
	}
}

func (dis *Discovery) Start() {
	fmt.Println("服务注册中心启动...")
	go func() {
		//获取服务器IP:port
		fmt.Println("获取服务注册中心TCP地址...")
		TcpAddr, err := net.ResolveTCPAddr(dis.IPVersion, fmt.Sprintf("%s:%d", dis.IP, dis.Port))
		if err != nil {
			fmt.Println("获取服务注册中心TCP地址出错 :", err)
			return
		}
		fmt.Println("获取服务注册中心TCP地址成功!")

		fmt.Println("服务注册中心获取listenner...")
		//network一般指TCP、HTTP的版本
		//监听端口
		listenner, err := net.ListenTCP(dis.IPVersion, TcpAddr)
		if err != nil {
			fmt.Println("服务注册中心获取listenner出错 :", err)
			return
		}
		fmt.Println("服务注册中心获取listenner成功!")

		//不停地获取链接conn
		for {
			fmt.Println("服务注册中心获取conn...")
			//TODO conn不能只做一次调用，要多次,加循环就可以了
			//TODO conn怎么判断已经断开？
			conn, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("服务注册中心获取conn出错:", err)
				return
			}
			fmt.Println("服务注册中心获取conn成功！")
			//疑问，为什么不需要传conn?闭包？
			go dis.DealConn(conn)
		}
	}()
}

func (dis *Discovery) Stop() {

}

func (dis *Discovery) Serve() {
	dis.Start()
	select {}
}

func (dis *Discovery) AddServer(req *DisRequest, addr string) error {
	_, err := dis.RedisClient.RPush(req.ServerName, addr).Result()
	return err
}

func (dis *Discovery) SearchServer(req *DisRequest) (string, error) {
	list := dis.RedisClient.LRange(req.ServerName, 0, -1).Val()
	if len(list) == 0 {
		return "", errors.New("没有这个服务！")
	}
	//负载均衡
	return list[0], nil
}

//TODO要读写分离
func (dis *Discovery) DealConn(conn *net.TCPConn) {
	//读出序列化的数据
	buf := make([]byte, 512)
	cnt, err := conn.Read(buf)
	if err != nil {
		fmt.Print("read msg from conn error :", err)
		return
	}
	//反序列化得到request
	//req的Msg没有初始化，却不报错？
	var req DisRequest
	err = json.Unmarshal(buf[:cnt], &req)
	if err != nil {
		fmt.Println("反序列化失败:", err)
		return
	}
	fmt.Println("反序列化成功！")

	res := NewDisResponse(nil, errors.New(""))

	//res会包含错误信息，所以这里不打印
	dis.DealRequest(conn, &req, res)
	fmt.Println("请求处理完毕，准备序列化后发送回去！")

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

func (dis *Discovery) DealRequest(conn *net.TCPConn, req *DisRequest, res *DisResponse) {
	flag := req.Flag
	//TODO 检查服务名是否合法
	//不能随便删除和注册服务，必须给权限
	//查找也是必须给权限
	switch flag {
	//注册服务
	case 1:
		addr := conn.RemoteAddr().String()
		err := dis.AddServer(req, addr)
		res.ErrMsg = err
	//查找服务
	case 2:
		addr, err := dis.SearchServer(req)
		serverAddr, _ := net.ResolveTCPAddr("tcp", addr)
		res.Addr = serverAddr
		res.ErrMsg = err
	//删除服务
	case 3:
		addr := conn.RemoteAddr().String()
		err := dis.DeleteServer(req, addr)
		res.ErrMsg = err
	default:
		res.ErrMsg = errors.New("没有这样的指令！")
	}
}
