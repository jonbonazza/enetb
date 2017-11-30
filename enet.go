package enet

// #cgo CFLAGS: -I"src/"
// #cgo LDFLAGS: lib/libenet.a
// #include "enet/enet.h"
import "C"

func Initialize() bool {
	return C.enet_initialize() == 0
}

func Deinitialize() {
	C.enet_deinitialize()
}
