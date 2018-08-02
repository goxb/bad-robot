package main

import (
	"context"
	"flag"
	"log"

	"protos"

	"github.com/smallnest/rpcx/client"
)

var (
	addr2 = flag.String("addr", "127.0.0.1:3345", "server address")
)

func main() {
	flag.Parse()

	d := client.NewPeer2PeerDiscovery("tcp@"+*addr2, "")
	xclient := client.NewXClient("MRobotRpcProtoMod",
		client.Failtry, client.RandomSelect, d, client.DefaultOption)
	defer xclient.Close()

	args := &protos.MRobotAllocRequest{
		Ptype:  18,
		Callid: "abcdef12345678",
	}

	reply := &protos.MRobotAllocResponse{}
	call, err := xclient.Go(context.Background(),
		"RpcCall_MRobotAllocRequest", args, reply, nil)
	if err != nil {
		log.Fatalf("failed to call: %v", err)
	}

	replyCall := <-call.Done
	if replyCall.Error != nil {
		log.Fatalf("failed to call: %v", replyCall.Error)
	} else {
		log.Printf("server response %+v", reply)
	}
}
