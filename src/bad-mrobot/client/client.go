package main

import (
	"context"
	"flag"
	"log"
	"time"

	"protos"

	"github.com/smallnest/rpcx/client"
)

var (
	addr = flag.String("addr", "127.0.0.1:3345", "server address")
)

func main() {
	flag.Parse()

	d := client.NewPeer2PeerDiscovery("tcp@"+*addr, "")
	xclient := client.NewXClient("MRobotRpcProtoMod",
		client.Failtry, client.RandomSelect, d, client.DefaultOption)
	defer xclient.Close()

	callid := "2af146f463616c6c007c75d8@120.78.226.202"
	args := &protos.MRobotAllocRequest{
		Ptype:  18,
		Callid: callid,
	}

	for {
		reply := &protos.MRobotAllocResponse{}
		err := xclient.Call(context.Background(),
			"RpcCall_MRobotAllocRequest", args, reply)
		if err != nil {
			log.Fatalf("failed to call: %v", err)
		}

		log.Printf("server response %+v", reply)

		req := &protos.MRobotSetRomoteRequest{
			Callid:    callid,
			RtpRobot:  "127.0.0.1:20000",
			RtpRemote: "127.0.0.1:40000",
		}

		rsp := &protos.MRobotSetRemoteResponse{}
		err = xclient.Call(context.Background(),
			"RpcCall_MRobotSetRomoteRequest", req, rsp)
		if err != nil {
			log.Fatalf("failed to call: %v", err)
		}

		log.Printf("server response %+v", rsp)
		time.Sleep(1e9)
	}
}
