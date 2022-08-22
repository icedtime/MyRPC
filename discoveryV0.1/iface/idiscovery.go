package iface

import "discovery/impl"

type Idiscovery interface {
	AddServer(req *impl.DisRequest, res *impl.DisResponse) error
	//感觉不需要删除功能，服务器就注册服务，客户端就查找服务
	//DeleteServer(req *impl.DisRequest, res *impl.DisResponse) error
	SearchServer(req *impl.DisRequest, res *impl.DisResponse) error
}
