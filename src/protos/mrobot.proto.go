package protos

const (
	SUCCESS = 0
)

type Response struct {
	Code   int32
	Errmsg string
}

type MediaChannel struct {
	Ip    string
	Port  int32
	Ptype int32
}

type MRobotAllocRequest struct {
	Ptype  int32
	Callid string
}

type MRobotAllocResponse struct {
	Response
	Callid string
	Amedia MediaChannel
}

type MRobotFreeRequest struct {
	Callid string
}

type MRobotFreeResponse struct {
	Response
}
