package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jonbonazza/enetb"
)

func main() {
	var mode string
	flag.StringVar(&mode, "mode", "server", "")
	flag.Parse()
	enet.Initialize()
	defer enet.Deinitialize()
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	if mode == "server" {
		server(ch)
	} else if mode == "client" {
		client(ch)
	}
}

func server(sig chan os.Signal) {
	host, err := enet.NewHost("127.0.0.1:8080", 32, 2)
	if err != nil {
		panic(err)
	}
	defer host.Destroy()
	go poll(host)
	<-sig
}

func client(sig chan os.Signal) {
	host, err := enet.NewHost("", 1, 2)
	if err != nil {
		panic(err)
	}
	go poll(host)
	rand.Seed(time.Now().Unix())
	id := rand.Int()
	peer, err := enet.Connect(host, "127.0.0.1:8080", 3, uint32(id))
	if err != nil {
		panic(err)
	}
	defer func() {
		peer.Disconnect()
		host.Flush()
		time.Sleep(100 * time.Millisecond)
		host.Destroy()
	}()
	time.Sleep(5 * time.Second)
	if err := peer.Write(1, []byte("test"), enet.PacketFlagReliable); err != nil {
		panic(err)
	}
	<-sig
}

func poll(host *enet.Host) {
	ch := host.Poll(1 * time.Second)
	for e := range ch {
		handleEvent(e)
	}
}

func handleEvent(e *enet.Event) {
	switch e.EventType {
	case enet.EventTypeConnect:
		fmt.Println("Connect", e.EventData)
	case enet.EventTypeDisconnect:
		fmt.Println("Disconnect", e.EventData)
	case enet.EventTypeReceive:
		fmt.Println("Receive")
		fmt.Println(string(e.Data), e.EventData)
	}
}
