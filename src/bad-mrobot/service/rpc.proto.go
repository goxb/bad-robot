// Package example defines datastructure and services.
package service

import (
	"context"
	"fmt"
	"protos"
	"runtime/debug"

	"bad-mrobot/session"

	"github.com/cihub/seelog"
)

type MRobotRpcProtoMod protos.MRobotRpcProtoMod

func (m *MRobotRpcProtoMod) RpcCall_MRobotFreeRequest(
	ctx context.Context,
	req *protos.MRobotFreeRequest,
	rsp *protos.MRobotFreeResponse) error {

	seelog.Infof("MRobotFreeRequest %+v", req)
	defer func() {
		if err := recover(); err != nil {
			rsp.Code = -1
			rsp.Errmsg = fmt.Sprintf("%s", err)
			seelog.Errorf("\n%s", debug.Stack())
		}

		seelog.Infof("MRobotFreeResponse %+v", rsp)
	}()

	rsp.Code = protos.SUCCESS
	rsp.Errmsg = protos.SUCCESS_MSG
	rsp.Callid = req.Callid

	session.DestroySession(req.Callid)
	return nil
}

func (m *MRobotRpcProtoMod) RpcCall_MRobotAllocRequest(
	ctx context.Context,
	req *protos.MRobotAllocRequest,
	rsp *protos.MRobotAllocResponse) error {

	seelog.Infof("MRobotAllocRequest %+v", req)
	defer func() {
		if err := recover(); err != nil {
			rsp.Code = -1
			rsp.Errmsg = fmt.Sprintf("%s", err)
			seelog.Errorf("\n%s", debug.Stack())
		}

		seelog.Infof("MRobotAllocResponse %+v", rsp)
	}()

	rsp.Code = protos.SUCCESS
	rsp.Errmsg = protos.SUCCESS_MSG
	rsp.Callid = req.Callid

	s := session.GetSession(req.Callid)
	if s != nil {
		goto success
	}

	s = session.CreateSession(req.Callid)
	if s == nil {
		rsp.Code = -1
		rsp.Errmsg = fmt.Sprintf(
			"create session error %s", req.Callid)
		seelog.Errorf(rsp.Errmsg)
		goto failure
	}

	s.SetPayloadType(req.Ptype)
	if err := s.Init(); err != nil {
		rsp.Errmsg = fmt.Sprintf(
			"init session error %s", req.Callid)
		seelog.Errorf(rsp.Errmsg)
		goto failure
	}

	s.Start()

success:
	rsp.Amedia.Ptype = req.Ptype
	rsp.Amedia.RtpRobot = fmt.Sprintf("%s:%d",
		s.RtpRobot2.IpAddr.String(), s.RtpRobot2.DataPort)
	rsp.Amedia.RtcpRobot = fmt.Sprintf("%s:%d",
		s.RtpRobot2.IpAddr.String(), s.RtpRobot2.CtrlPort)
	seelog.Debugf("allocate session %p", s)
	return nil

failure:
	session.DestroySession(req.Callid)
	return nil
}

func (m *MRobotRpcProtoMod) RpcCall_MRobotSetRomoteRequest(
	ctx context.Context,
	req *protos.MRobotSetRomoteRequest,
	rsp *protos.MRobotSetRemoteResponse) error {

	seelog.Infof("MRobotSetRomoteRequest %+v", req)
	defer func() {
		if err := recover(); err != nil {
			rsp.Code = -1
			rsp.Errmsg = fmt.Sprintf("%s", err)
			seelog.Errorf("\n%s", debug.Stack())
		}

		seelog.Infof("MRobotSetRemoteResponse %+v", rsp)
	}()

	rsp.Code = protos.SUCCESS
	rsp.Errmsg = protos.SUCCESS_MSG
	rsp.Callid = req.Callid
	s := session.GetSession(req.Callid)
	if s == nil {
		rsp.Code = -1
		rsp.Errmsg = fmt.Sprintf(
			"get session error %s", req.Callid)
		seelog.Errorf(rsp.Errmsg)
		return nil
	}

	s.AddRemote(req.RtpRobot, req.RtpRemote)
	seelog.Debugf("set remote session %p", s)
	return nil
}
