package enet

// #cgo CFLAGS: -I"src/"
// #cgo LDFLAGS: lib/libenet.a
// #include "enet/enetb.h"
import "C"
import (
	"errors"
	"net"
	"strconv"
	"unsafe"
)

func newAddress(addr string) (*C.ENetAddress, error) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}
	chost := C.CString(host)
	defer C.free(unsafe.Pointer(chost))
	p, err := strconv.Atoi(port)
	if err != nil {
		return nil, err
	}
	enetAddr := C.create_enet_address(chost, C.enet_uint16(p))
	if enetAddr == nil {
		return nil, errors.New("failed to create enet address")
	}
	return enetAddr, nil
}

func getAddress(addr C.ENetAddress) string {
	host := getHost(addr)
	port := strconv.Itoa(int(addr.port))
	return host + ":" + port
}

func getHost(addr C.ENetAddress) string {
	size := 16
	chost := (*C.char)(C.malloc(C.size_t(size)))
	ret := C.enet_address_get_host_ip(&addr, chost, C.size_t(size))
	if ret != 0 {
		return ""
	}
	defer C.free(unsafe.Pointer(chost))
	return C.GoStringN(chost, C.int(size))
}
