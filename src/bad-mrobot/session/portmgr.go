package session

var gUdpPort int = 20000

func GetZone() string {
	return ""
}

func GetFreeUdpPort() int {
	gUdpPort += 2
	return gUdpPort
}

func GetWanIpAddress() string {
	return "127.0.0.1"
}
