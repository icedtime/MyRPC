package impl

import "net"

type DisRequest struct {
	ServerName string
	Flag       int
}

func NewDisRequest(serverName string, flag int) *DisRequest {
	return &DisRequest{
		ServerName: serverName,
		Flag:       flag,
	}
}

type DisResponse struct {
	ErrMsg error
	Addr   *net.TCPAddr
}

func NewDisResponse(addr *net.TCPAddr, err error) *DisResponse {
	return &DisResponse{
		Addr:   addr,
		ErrMsg: err,
	}
}
