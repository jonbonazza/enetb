package enet

// #cgo CFLAGS: -I"src/"
// #cgo LDFLAGS: lib/libenet.a
// #include "enet/enet.h"
import "C"
import "unsafe"

const (
	EventTypeNone = iota
	EventTypeConnect
	EventTypeDisconnect
	EventTypeReceive
)

type Event struct {
	EventType int
	Peer      *Peer
	ChannelID uint8
	EventData uint32
	Flags     int
	Data      []byte
}

func newEvent(e C.ENetEvent) *Event {
	var data []byte
	if e.packet != nil {
		data = C.GoBytes(unsafe.Pointer(e.packet.data), C.int(e.packet.dataLength))
	}
	return &Event{
		EventType: int(e._type),
		Peer:      &Peer{e.peer},
		ChannelID: uint8(e.channelID),
		EventData: uint32(e.data),
		Data:      data,
	}
}
