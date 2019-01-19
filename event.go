package enet

// #cgo CFLAGS: -I"src/"
// #cgo LDFLAGS: lib/libenet.a
// #include "enet/enet.h"
import "C"
import "unsafe"

const (
	// EventTypeNone is returned if no event occurred within the specified time limit.
	EventTypeNone = iota
	// EventTypeConnect is returned when either a client host has connected to the
	// server host or when an attempt to establish a connection with a foreign host has succeeded.
	EventTypeConnect
	// EventTypeDisconnect is returned when a connected peer has either explicitly disconnected or timed out.
	EventTypeDisconnect
	//EventTypeReceive is returned when a packet is received from a connected peer.
	EventTypeReceive
)

// Event represents a single event from the ENet subsystem.
type Event struct {
	// The type of the event. Will be one of EventTypeNone, EventTypeConnect,
	// EventTypeDisconnect, or EventTypeReceive.
	EventType int
	// Contains the information for the remote peer. In the case of EventTypeConnect,
	// Peer will be the newly connected peer. In the case of EventTypeReceive,
	// Peer will be the peer that the packet was sent from. In the case of EventTypDisconnect,
	// Peer will be the peer that disconnected.
	Peer *Peer
	// ChannelID is only used for EventTypeReceive and is the id of the channel that the packet
	// was received from.
	ChannelID uint8
	// EventData is arbitrary data sent a long with the event.
	EventData uint32
	// Contains the raw packet bytes.
	Data []byte
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
