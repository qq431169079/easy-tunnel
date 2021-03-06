package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
)

var (
	remoteHost     *string
	remotePort     *int
	remoteBindHost *string
	remoteBindPort *int
	localHost      *string
	localPort      *int
	tunnelType	   *string
)

func init() {
	remoteHost = flag.String("h", "127.0.0.1", "远程服务器通信ip")
	remotePort = flag.Int("p", 9960, "远程服务器通信端口")
	remoteBindHost = flag.String("bh", "0.0.0.0", "开启映射后绑定ip")
	remoteBindPort = flag.Int("bp", -1, "远程开启的映射端口,必填")
	localHost = flag.String("fh", "127.0.0.1", "转发目标ip")
	localPort = flag.Int("fp", -1, "转发目标端口,必填")
	tunnelType = flag.String("protocol", "tcp", "映射通信协议")

	log.SetPrefix("TRACE: ")
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)
}

func main() {
	flag.Parse()

	if *localPort == -1 {
		panic("fp参数必填")
	}
	if *remoteBindPort == -1 {
		panic("bp参数必填")
	}

	bridgeClient := NewBridgeClient()
	success, err := bridgeClient.Connect(*remoteHost, *remotePort)
	if !success {
		panic(err)
	}
	msg, e := bridgeClient.OpenTunnel(*remoteBindHost, *remoteBindPort, *localHost, *localPort, *tunnelType)
	if e != nil {
		panic(msg)
	}

	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for _ = range signalChan {
			log.Println("exit bridge server")
			bridgeClient.Close()
			cleanupDone <- true
		}
	}()
	<-cleanupDone
}
