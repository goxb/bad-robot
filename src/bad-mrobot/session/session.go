package session

import (
	"net"

	"github.com/wernerd/GoRTP/src/net/rtp"
)

type RobotSession struct {
	Ptype  int32
	Callid string
	rtport *rtp.Session
}

type SessionMgr map[string]*RobotSession

func (this *RobotSession) Init() error {
	zone := GetZone()
	port := GetFreeUdpPort()
	wanip := GetWanIpAddress()

	local, err = net.ResolveIPAddr("ip", wanip)
	if err != nil {
		return err
	}

	// Create a UDP transport with "local" address and use this for a "local" RTP session
	// The RTP session uses the transport to receive and send RTP packets to the remote peer.
	tpLocal, err := rtp.NewTransportUDP(local, port, GetZone())
	if err != nil {
		return err
	}

	// TransportUDP implements TransportWrite and TransportRecv interfaces thus
	// use it to initialize the Session for both interfaces.
	rsLocal := rtp.NewSession(tpLocal, tpLocal)

	// Create a media stream.
	// The SSRC identifies the stream. Each stream has its own sequence number and other
	// context. A RTP session can have several RTP stream for example to send several
	// streams of the same media.
	//

	Idx, err := rsLocal.NewSsrcStreamOut(
		&rtp.Address{wanip, port, port + 1, GetZone()}, 1020304, 4711)
	if err != nil {
		return err
	}

	rsLocal.SsrcStreamOutForIndex(Idx).SetPayloadType(0)
	this.rtport = rsLocal
	return nil
}

func (this *RobotSession) AddRemote(ip string, port int) error {
	// Add address of a remote peer (participant)
	this.rtport.AddRemote(&rtp.Address{ip, port, port + 1, GetZone()})
}
