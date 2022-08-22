package discoveryv02

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
)

type RedisDiscovery struct {
	hearBeat uint          //sec，监听redis是否工作正常
	pool     *redis.Pool   //redis连接池
	close    chan struct{} //判断当前服务注册中心是否需要关闭

	ser *Server
	mu  sync.Mutex
}

func NewRedisDiscovery(hearBeat uint, addr string, auth *string) *RedisDiscovery {
	dis := &RedisDiscovery{
		hearBeat: 120, //默认为3s，但是如果给的大于3s，就用大的那个
		close:    make(chan struct{}),
	}

	if hearBeat >= 3 {
		dis.hearBeat = hearBeat
	}
	pool := &redis.Pool{
		MaxIdle:     10,                // 最大空闲连接数
		MaxActive:   10,                // 最大连接数
		IdleTimeout: 300 * time.Second, // 超时回收
		Dial: func() (redis.Conn, error) {
			dial, err := redis.Dial("tcp", addr)
			if err != nil {
				fmt.Println("Redis Pool err:", err)
				return nil, err
			}
			//访问认证
			if auth != nil {
				dial.Do("AUTH", *auth)
			}
			return dial, nil
		},

		//time.Time本质就是一个时间，time.Time有3个字段：
		//sec表示公元1年1月1日00:00:00UTC到要表示的整数秒数
		//nsec则表示余下的纳米描述
		//loc就是时区，就是偏移值
		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			//time.Since(t)表示过了t多少时间
			if time.Since(t) < time.Minute {
				return nil
			}
			//不然就ping一次redis，看当前链接是否可用
			_, err := conn.Do("PING")
			return err
		},
	}
	dis.pool = pool

	return dis
}

func (r *RedisDiscovery) Test() {
	key, value := "PangYin", "good luck！"
	rds := r.pool.Get()
	defer rds.Close()

	_, err := rds.Do("set", key, value)
	if err != nil {
		fmt.Println(err)
		return
	}
}

//查找服务功能
//根据服务名查找，提供服务的应该是一组服务器
func (r *RedisDiscovery) Discovery(serverName string) ([]*Server, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	//存在redis的形式是 key:todaaa/serverName/ID,value:server{}
	//但redis支持查找key为todaaa/serverName/*的查找方式
	path := fmt.Sprintf("todaaa/%s/*", serverName)
	rds := r.pool.Get()
	defer rds.Close()

	sers := make([]*Server, 0)

	values, err := redis.Strings(rds.Do("keys", path))
	if err != nil {
		return nil, err
	}

	//values应该是一组key的集合，就是符合上面todaaa/serverName/*的一组key
	for _, v := range values {
		//用这些key查找
		buf, err := redis.Bytes(rds.Do("get", v))
		if err != nil {
			return nil, err
		}

		ser := Server{}
		err = json.Unmarshal(buf, &ser)
		if err != nil {
			fmt.Println(err)
			continue
		}

		sers = append(sers, &ser)
	}
	return sers, nil
}

//注册功能
//一个服务包含的字段有：服务器ID，服务名，服务器addr，权重，允许的最大负载和当前负载
//后三项可以用来做负载均衡
//对于一项服务怎么存在redis？
//按照这样的格式：仓库名/服务名/服务器ID，目前来说:todaaa/serverName/ID
func (r *RedisDiscovery) Registry(serverName, addr string, weights float64, maximumLoad, CurrentLoad int64, serverID *string) error {
	//判断id
	if serverID == nil {
		//id, err := utils.DistributedID()
		id := "123456789"
		// if err != nil {
		// 	return err
		// }
		serverID = &id
	}

	r.ser = &Server{
		ServerName: serverName,
		Addr:       addr,
		ID:         *serverID,
		Weights:    weights,
		//Protocol:    protocol,
		MaximumLoad: maximumLoad,
	}

	//todaaa/serverName/serverID
	path := r.getRedisPath(serverName, *serverID)
	err := r.registry(path)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-r.close:
				return
			case <-time.After(time.Second * time.Duration(3)):
				err = r.registry(path)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}()

	return nil
}

func (r *RedisDiscovery) registry(path string) error {
	//从连接池获取一个连接
	rds := r.pool.Get()
	defer rds.Close()

	//使用json序列化为[]byte
	buf, err := json.Marshal(r.ser)
	if err != nil {
		return err
	}

	//Redis Setex 命令为指定的 key 设置值及其过期时间。如果 key 已经存在， SETEX 命令将会替换旧的值。
	_, err = rds.Do("setex", path, r.hearBeat, buf)
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisDiscovery) getRedisPath(serverName string, serverID string) string {
	return fmt.Sprintf("todaaa/%s/%s", serverName, serverID)
}

// func (r *RedisDiscovery) UnRegistry(serName string, serID string) error {
// 	close(r.close)
// 	return nil
// }

// func (r *RedisDiscovery) Add(load int64) {
// 	atomic.AddInt64(&r.ser.CurrentLoad, load)
// }
// func (r *RedisDiscovery) Less(load int64) {
// 	atomic.AddInt64(&r.ser.CurrentLoad, -load)
// }
//限流用的
// func (r *RedisDiscovery) Limit() bool
