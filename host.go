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

// Host is a wrapper around a *C.ENetHost and represents a network host.
// Depending on its usage, it can be either a server host or a client host.
//
// See the ENet documentation for more information.
type Host struct {
	chost *C.ENetHost
}

// Destroy destroys the host and frees all of its resources. This MUST be called
// Before the Host is garbage collected by the Go Runtime, or memory and other resources
// will be leaked.
func (h *Host) Destroy() {
	C.enet_host_destroy(h.chost)
}

// Addr returns the host's address as a string.
func (h *Host) Addr() string {
	return getAddress(h.chost.address)
}

// Poll begins polling the Host for peer events. If the provided timeout is reached before an event
// is received, an event of type EventTypeNone is received. Poll is non-blocking and a channel is returned that can be used to listen for events.
// This channel will be closed once polling has stopped. If timeout is 0, polling will stop and the channel will be closed immediately if there are
// no events to dispatch.
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
			if ret == 0 {
				// Timeout waiting for event, continue.
				if timeout != 0 {
					continue
				}
				// 0 timeout and no events. Stop polling.
				return
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

// Flush flushes any currnetly queued packets without emitting events. This is useful for
// gracefully shutting down when you don't care about any packages received after the shutdown intent.
func (h *Host) Flush() {
	C.enet_host_flush(h.chost)
}

// Broadcast sents a packet with the provided data to all currently connected peers on the channel with channelID.
// Flags are the flags that are used when creating the ENet packet. See the ENet documentation for more information
// on the various flags and what they do.
func (h *Host) Broadcast(channelID uint8, data []byte, flags int) {
	packet := C.enet_packet_create(C.CBytes(data), C.size_t(len(data)), C.enet_uint32(flags))
	C.enet_host_broadcast(h.chost, C.enet_uint8(channelID), packet)
}

// NewHost returns a new host. The host will listen on addr and can be connected to up to maxClients peers at once.
// The number of channels can also be specified.
//
// An error is returned if the host could not be created.
func NewHost(addr string, maxClients, numChannels uint32) (*Host, error) {
	var enetAddr *C.ENetAddress
	if addr != "" {
		var err error
		enetAddr, err = newAddress(addr)
		if err != nil {
			return nil, err
		}
	}
	enetHost := C.enet_host_create(enetAddr, C.size_t(maxClients), C.size_t(numChannels), 0, 0)
	if enetHost == nil {
		if enetAddr != nil {
			C.free(unsafe.Pointer(enetAddr))
		}
		return nil, errors.New("failed to create Enet host")
	}
	return &Host{
		chost: enetHost,
	}, nil
}
