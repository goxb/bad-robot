package service

import (
	"github.com/smallnest/rpcx/server"
)

type RpcService struct {
	address string
	rpcxsrv *server.Server
}

var gRpcService *RpcService

func (s *RpcService) reg(rcvr interface{}) {
	s.rpcxsrv.Register(rcvr, "")
}

func (s *RpcService) run() error {
	return s.rpcxsrv.Serve("tcp", s.address)
}

func (s *RpcService) register() {
	s.reg(new(MRobotRpcProtoMod))
}

func RunRpcService(address string) error {
	gRpcService = &RpcService{
		address: address,
		rpcxsrv: server.NewServer(),
	}

	gRpcService.register()
	return gRpcService.run()
}
