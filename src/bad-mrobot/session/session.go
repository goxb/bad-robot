package session

import (
	"io/ioutil"
	"net"
	"sync"
	"time"

	"github.com/cihub/seelog"
	"github.com/wernerd/GoRTP/src/net/rtp"
)

type RobotSession struct {
	id        string
	lock      sync.RWMutex
	session   *rtp.Session
	stopRecv  chan bool
	stopCtrl  chan bool
	remoteMap map[string]string

	ctime int64
	ptype int32

	RtpRobot2 rtp.Address
}

var eventNamesNew = []string{"NewStreamData", "NewStreamCtrl"}
var eventNamesRtcp = []string{"SR", "RR", "SDES", "BYE"}

func (this *RobotSession) Free() {
	close(this.stopRecv)
	close(this.stopCtrl)
	this.session.CloseSession()
}

func (this *RobotSession) Id() string {
	return this.id
}

func (this *RobotSession) Init() error {
	port := GetFreeUdpPort()
	wanip := GetWanIpAddress()

	addr, err := net.ResolveIPAddr("ip", wanip)
	if err != nil {
		seelog.Errorf("resolve ip addr error %s %s", err, this.Id())
		return err
	}

	// Create a UDP transport with "addr" address and use this for a "addr" RTP session
	// The RTP session uses the transport to receive and send RTP packets to the remote peer.
	tp, err := rtp.NewTransportUDP(addr, port, GetZone())
	if err != nil {
		seelog.Errorf("new transport udp error %s %s", err, this.Id())
		return err
	}

	// TransportUDP implements TransportWrite and TransportRecv interfaces thus
	// use it to initialize the Session for both interfaces.
	session := rtp.NewSession(tp, tp)

	// Create a media stream.
	// The SSRC identifies the stream. Each stream has its own sequence number and other
	// context. A RTP session can have several RTP stream for example to send several
	// streams of the same media.
	//

	this.RtpRobot2 = rtp.Address{addr.IP, port, port + 1, GetZone()}
	idx, err2 := session.NewSsrcStreamOut(&this.RtpRobot2, 1020304, 4711)
	if err2 != "" {
		seelog.Errorf("new ssrc stream out %s %s", err2, this.Id())
		return err2
	}

	session.SsrcStreamOutForIndex(idx).SetPayloadType(byte(this.ptype))
	session.StartSession()
	this.session = session
	return nil
}

func (this *RobotSession) AddRemote(robot, remote string) error {
	// Add address of a remote peer (participant)
	addr, err := net.ResolveUDPAddr("udp", remote)
	if err != nil {
		return err
	}

	this.lock.Lock()
	defer this.lock.Unlock()

	if _, ok := this.remoteMap[remote]; !ok {
		this.session.AddRemote(
			&rtp.Address{
				IpAddr:   addr.IP,
				DataPort: addr.Port,
				CtrlPort: addr.Port + 1,
				Zone:     GetZone(),
			})

		this.remoteMap[remote] = robot
		seelog.Infof("add %s remote %s", robot, remote)
	}

	seelog.Debugf("remoteMap size %d", len(this.remoteMap))
	return nil
}

func (this *RobotSession) SetPayloadType(t int32) {
	this.ptype = t
}

func LoadG729() []byte {
	ring := "./timelimit.g729"
	fileBuf, err := ioutil.ReadFile(ring)
	if err != nil {
		seelog.Error("ReadFile ", ring, " error ", err)
		return nil
	}

	return fileBuf
}

// Create a RTP packet suitable for standard stream (index 0) with a payload length of 160 bytes
// The method initializes the RTP packet with SSRC, sequence number, and RTP version number.
// If the payload type was set with the RTP stream then the payload type is also set in
// the RTP packet
func (this *RobotSession) SendData() {
	stamp := uint32(0)
	fileBuf := LoadG729()
	var g729_payload byte = 18
	var g729_frame_size int = 20

	ringlen := len(fileBuf)
	count := ringlen / g729_frame_size
	if ringlen%g729_frame_size != 0 {
		count++
	}

	for k := 0; k < 10; k++ {
		for i := 0; i < count; i++ {
			beg := i * g729_frame_size
			end := beg + g729_frame_size
			if end > ringlen {
				end = ringlen
			}

			paload := make([]byte, end-beg)
			copy(paload, fileBuf[beg:end])

			rp := this.session.NewDataPacket(stamp)
			rp.SetPayload(paload[:])
			rp.SetPayloadType(g729_payload)
			this.session.WriteData(rp)
			rp.FreePacket()

			stamp += 160
			time.Sleep(20e6)
		}

		seelog.Infof("play ring finish %d %s", k, this.Id())
	}

	seelog.Infof("play ring all finish %s", this.Id())
}

func (this *RobotSession) ReceiveData() {
	// Create and store the data receive channel.
	dataReceiver := this.session.CreateDataReceiveChan()
	var cnt int

	for {
		select {
		case rp := <-dataReceiver: // just get a packet - maybe we add some tests later
			if (cnt % 50) == 0 {
				seelog.Infof("Remote receiver got %d packets", cnt)
			}
			cnt++
			rp.FreePacket()
		case <-this.stopRecv:
			seelog.Warnf("stop data rcv %s", this.Id())
			return
		}
	}
}

func (this *RobotSession) ReceiveCtrl() {
	// Create and store the control event channel.
	ctrlReceiver := this.session.CreateCtrlEventChan()
	for {
		select {
		case evSlice := <-ctrlReceiver: // get an event
			seelog.Infof("length of event slice:", len(evSlice))
			for _, event := range evSlice {
				if event != nil {
					var eventName string
					if event.EventType < 200 {
						eventName = eventNamesNew[event.EventType]
					} else {
						eventName = eventNamesRtcp[event.EventType-200]
					}
					seelog.Infof("received ctrl event, type: %s, ssrc: %d, %s",
						eventName, event.Ssrc, event.Reason)
				} else {
					seelog.Infof("unexpected nil event")
				}
			}
		case <-this.stopCtrl:
			seelog.Warnf("stop ctrl rcv %s", this.Id())
			return
		}
	}
}

func (this *RobotSession) Start() {
	go this.ReceiveData()
	go this.ReceiveCtrl()
	go this.SendData()
}
