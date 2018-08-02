package session

import "bad-mrobot/config"

var gUdpPort int = 20000

func GetZone() string {
	return ""
}

func GetFreeUdpPort() int {
	gUdpPort += 2
	return gUdpPort
}

func GetWanIpAddress() string {
	return config.RtpIpAddr
}
