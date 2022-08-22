package discoveryv02

type IDdiscovery interface {
	Discovery(serName string) ([]*Server, error)
	Registry(serName, addr string, weights float64 /*, protocol transport.Protocol*/, maximumLoad int64, serID *string) error
	UnRegistry(serName string, serID string) error
	//TODO 限流
	//Limit() bool // 是否开启限流 Whether to limit the flow
	Add(load int64)
	Less(load int64)
}

type Server struct {
	ID         string `json:"id"`          // 体用服务的服务器ID，由算法生成
	ServerName string `json:"server_name"` // 服务名称
	Addr       string `json:"addr"`        // 地址
	//Protocol    transport.Protocol `json:"protocol"`     // 协议

	//负载均衡用的
	Weights     float64 `json:"weights"`      // 权重
	MaximumLoad int64   `json:"maximum_load"` // 最大负载
	CurrentLoad int64   `json:"current_load"` // 当前负载
}
