package enet

// #cgo CFLAGS: -I"src/"
// #cgo LDFLAGS: lib/libenet.a
// #include "enet/enet.h"
import "C"
import (
	"errors"
	"fmt"
	"unsafe"
)

const (
	PeerStateDisconnected = iota
	PeerStateConnecting
	PeerStateAcknowledgingConnect
	PeerStateConnectionPending
	PeerStateConnectionSucceeded
	PeerStateConnected
	PeerStateDisconnectLater
	PeerStateDisconnecting
	PeerStateAcknowledgingDisconnect
	PeerStateZombie

	PacketFlagReliable           = 1 << 0
	PacketFlagUnsequenced        = 1 << 1
	PacketFlagNoAllocate         = 1 << 2
	PacketFlagUnreliableFragment = 1 << 3
)

type Peer struct {
	cpeer *C.ENetPeer
}

func (p *Peer) ID() uint32 {
	return uint32(p.cpeer.connectID)
}

func (p *Peer) Addr() string {
	return getAddress(p.cpeer.address)
}

func (p *Peer) State() int {
	return int(p.cpeer.state)
}

func (p *Peer) Write(channelID uint8, data []byte, flags int) error {
	packet := C.enet_packet_create(C.CBytes(data), C.size_t(len(data)), C.enet_uint32(flags))
	if ret := C.enet_peer_send(p.cpeer, C.enet_uint8(channelID), packet); ret != 0 {
		return fmt.Errorf("failed to send data over channel %d", channelID)
	}
	return nil
}

func (p *Peer) Disconnect(data uint32) {
	C.enet_peer_disconnect_now(p.cpeer, C.enet_uint32(data))
}

func Connect(host *Host, addr string, channelCount int, id uint32) (*Peer, error) {
	enetAddr, err := newAddress(addr)
	if err != nil {
		return nil, err
	}
	cpeer := C.enet_host_connect(host.chost, enetAddr, C.size_t(channelCount), C.enet_uint32(id))
	if cpeer == nil {
		C.free(unsafe.Pointer(enetAddr))
		return nil, errors.New("failed to create peer")
	}
	return &Peer{cpeer}, nil
}
