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

	args := &protos.MRobotAllocRequest{
		Ptype:  18,
		Callid: "abcdef12345678",
	}

	for {
		reply := &protos.MRobotAllocResponse{}
		err := xclient.Call(context.Background(),
			"RpcCall_MRobotAllocRequest", args, reply)
		if err != nil {
			log.Fatalf("failed to call: %v", err)
		}

		log.Printf("server response %+v", reply)
		time.Sleep(1e9)
	}
}
