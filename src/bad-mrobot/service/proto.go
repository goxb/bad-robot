// Package example defines datastructure and services.
package service

import (
	"context"
	"fmt"
	"protos"
)

type MRobotRpcProtoMod protos.MRobotRpcProtoMod

func (m *MRobotRpcProtoMod) RpcCall_MRobotAllocRequest(
	ctx context.Context,
	req *protos.MRobotAllocRequest,
	rsp *protos.MRobotAllocResponse) error {

	rsp.Code = protos.SUCCESS
	rsp.Errmsg = "success"
	rsp.Callid = req.Callid

	rsp.Amedia.Ip = "127.0.0.1"
	rsp.Amedia.Port = 3345
	rsp.Amedia.Ptype = req.Ptype

	fmt.Printf("client request: %+v\n", req)
	return nil
}
