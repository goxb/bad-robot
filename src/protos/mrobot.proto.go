package protos

const (
	SUCCESS     = 0
	SUCCESS_MSG = "success"
)

type Response struct {
	Code   int32
	Errmsg string
}

type MediaChannel struct {
	Ptype     int32
	RtpRobot  string
	RtcpRobot string
}

type MRobotFreeRequest struct {
	Callid string
}

type MRobotFreeResponse struct {
	Response
	Callid string
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

type MRobotSetRomoteRequest struct {
	Callid    string
	RtpRobot  string
	RtpRemote string
}

type MRobotSetRemoteResponse struct {
	Response
	Callid string
}
