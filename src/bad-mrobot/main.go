package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cihub/seelog"

	"bad-mrobot/config"
	"bad-mrobot/log"
	"bad-mrobot/service"
)

var (
	help bool

	beg int
	end int
	rtp string
	srv string
	cfg string

	logfile  string
	dbglevel string
)

func init() {
	flag.BoolVar(&help, "help", false, "this help")

	flag.IntVar(&beg, "beg", 20000, "beg udp port")
	flag.IntVar(&end, "end", 65535, "end udp port")
	flag.StringVar(&rtp, "rtp", "127.0.0.1", "rtp wan address")
	flag.StringVar(&srv, "srv", "127.0.0.1:3345", "tcp service address")
	flag.StringVar(&cfg, "cfg", "etc/bad-mrobot.conf", "set config file")
	flag.StringVar(&logfile, "log", "log/bad-mrobot.log", "set log file")
	flag.StringVar(&dbglevel, "dbg", "debug", "debug level: debug info warn error")

	// 改变默认的 Usage
	flag.Usage = usage
}

func usage() {
	command := filepath.Base(os.Args[0])
	fmt.Fprintf(os.Stderr, `%s
Usage: %s [-srv address] [-beg port] [-end port] [-dbg debug] [-log file]

Options:
`, command, command, command)
	flag.PrintDefaults()
}

func main() {
	flag.Parse()

	if help {
		flag.Usage()
		return
	}

	log.InitLog(
		dbglevel, logfile)

	defer seelog.Flush()

	config.RtpIpAddr = rtp
	err := service.RunRpcService(srv)
	seelog.Infof("RunRpcService %s", err)
}
