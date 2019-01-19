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

// The various states that a peer can be in at any given time.
// See the ENet documentation for more info.
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
)

// The various flags that can be used when creating a packet.
// See the ENet documentation for more info.
const (
	PacketFlagReliable           = 1 << 0
	PacketFlagUnsequenced        = 1 << 1
	PacketFlagNoAllocate         = 1 << 2
	PacketFlagUnreliableFragment = 1 << 3
)

// Peer is a wrapper around a *C.ENetPeer and represents a remote peer.
type Peer struct {
	cpeer *C.ENetPeer
}

// ID returns the peer's unique ID.
func (p *Peer) ID() uint32 {
	return uint32(p.cpeer.connectID)
}

// Addr is the remote address of the peer.
func (p *Peer) Addr() string {
	return getAddress(p.cpeer.address)
}

// State is the current state of the peer. It can be any one of the PeerState*
// constants defined above.
func (p *Peer) State() int {
	return int(p.cpeer.state)
}

// Write will asynchronously write data to the specified channel. Flags is a
// bitmask consisting of any number of the PacketFlag* constants defined above.
//
// An error is returned if the data could not be sent over the desired channel.
func (p *Peer) Write(channelID uint8, data []byte, flags int) error {
	packet := C.enet_packet_create(C.CBytes(data), C.size_t(len(data)), C.enet_uint32(flags))
	if ret := C.enet_peer_send(p.cpeer, C.enet_uint8(channelID), packet); ret != 0 {
		return fmt.Errorf("failed to send data over channel %d", channelID)
	}
	return nil
}

// Disconnect severs the peer's connection to the remote host.
// This MUST be called before the Peer is garbage collected, or
// socket and other types of leaks will occur.
func (p *Peer) Disconnect() {
	C.enet_peer_disconnect_now(p.cpeer, 0)
}

// Connect creates a connection to a remote host. ClientHost must
// be a client Host (no address) and addr is the host's address.
// Channel count is the number of channels that will be created
// for sending packets. Data is arbitrary data that can be associated with the Peer.
func Connect(clientHost *Host, addr string, channelCount int, data uint32) (*Peer, error) {
	enetAddr, err := newAddress(addr)
	if err != nil {
		return nil, err
	}
	cpeer := C.enet_host_connect(clientHost.chost, enetAddr, C.size_t(channelCount), C.enet_uint32(data))
	if cpeer == nil {
		C.free(unsafe.Pointer(enetAddr))
		return nil, errors.New("failed to create peer")
	}
	return &Peer{cpeer}, nil
}
