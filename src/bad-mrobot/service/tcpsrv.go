package service

import (
	"net"

	"github.com/cihub/seelog"
)

func Listen(address string) {
	l, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}

	for {
		accept(l)
	}
}

func accept(l net.Listener) {
	conn, err := l.Accept()
	if err != nil {
		seelog.Error(err)
		return
	}

	go readMessage(conn)
}

func readMessage(conn net.Conn) {
	defer conn.Close()
	/*
		buf := make([]byte, 65535)
		for {
			cnt, err := conn.Read(buf)
			if err != nil {
				seelog.Error(err)
				break
			}

			stReceive := &protos.UserInfo{}
			pData := buf[:cnt]

			//protobuf解码
			err = proto.Unmarshal(pData, stReceive)
			if err != nil {
				seelog.Error(err)
				break
			}
		}
	*/
}
