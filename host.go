package enet

// #cgo CFLAGS: -I"src/"
// #cgo LDFLAGS: lib/libenet.a
// #include "enet/enet.h"
import "C"
import (
	"errors"
	"time"
	"unsafe"
)

const (
	EnetEventTypeNone = iota
	EnetEventTypeConnect
	EnetEventTypeDisconnect
	EnetEventTypeReceive

	EnetHostAny = ""
)

var (
	ErrFailedToCreateHost = errors.New("failed to create Enet host")
)

type EnetEventType int

type PeerConnectedHandler func(host *Host, addr string, id uint32)
type PeerDisconnectedHandler func(host *Host, id uint32)

// chanid == 0xff unreliable data
type ReceiveHandler func(host *Host, fromAddr string, fromID uint32, chanid uint8, payload []byte)

type Host struct {
	chost *C.ENetHost
}

func (h *Host) Destroy() {
	C.enet_host_destroy(h.chost)
}

func (h *Host) Addr() string {
	return getAddress(h.chost.address)
}

func (h *Host) Poll(timeout time.Duration) <-chan *Event {
	ch := make(chan *Event, 256)
	go func() {
		defer close(ch)
		var event C.ENetEvent
		for {
			ret := C.enet_host_service(h.chost, &event, C.enet_uint32(timeout/time.Millisecond))
			// Error
			if ret < 0 {
				return
			}
			// Timedout waiting for event. Try again.
			if ret == 0 {
				continue
			}
			ch <- newEvent(event)
			// If we have a packet, we need to free it.
			if event.packet != nil {
				C.free(unsafe.Pointer(event.packet))
			}
			// If this is a disconnect type event, we need to reset the client information
			if int(event._type) == EventTypeDisconnect {
				event.peer.data = nil
			}
		}
	}()
	return ch
}

func (h *Host) Flush() {
	C.enet_host_flush(h.chost)
}

func (h *Host) Broadcast(channelID uint8, data []byte, flags int) {
	packet := C.enet_packet_create(C.CBytes(data), C.size_t(len(data)), C.enet_uint32(flags))
	C.enet_host_broadcast(h.chost, C.enet_uint8(channelID), packet)
}

func NewHost(addr string, maxClients uint32) (*Host, error) {
	var enetAddr *C.ENetAddress
	if addr != "" {
		var err error
		enetAddr, err = newAddress(addr)
		if err != nil {
			return nil, err
		}
	}
	enetHost := C.enet_host_create(enetAddr, C.size_t(maxClients), 2, 0, 0)
	if enetHost == nil {
		if enetAddr != nil {
			C.free(unsafe.Pointer(enetAddr))
		}
		return nil, ErrFailedToCreateHost
	}
	return &Host{
		chost: enetHost,
	}, nil
}
