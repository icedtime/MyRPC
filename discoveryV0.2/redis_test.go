package discoveryv02

import (
	"fmt"
	"testing"
)

func TestRedis(t *testing.T) {
	discovery := NewRedisDiscovery(3, "192.168.75.102:6379", nil)

	// err := discovery.Registry("Ser", "127.0.0.1:8654", 10, 512, 100, nil)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	//discovery.Test()
	err := discovery.Registry("PangYin", "127.0.0.1:8888", float64(0), int64(512), int64(100), nil)
	if err != nil {
		fmt.Println("111", err)
	}
	select {}

	// for {
	// 	servers, err := discovery.Discovery("Ser")
	// 	if err != nil {
	// 		log.Fatalln(err)
	// 		return
	// 	}
	// 	fmt.Println(len(servers))
	// 	time.Sleep(time.Second)
	// }
}
