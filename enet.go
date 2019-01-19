package enet

// #cgo CFLAGS: -I"src/"
// #cgo LDFLAGS: lib/libenet.a
// #include "enet/enet.h"
import "C"

// Initialize initializes the enet subsystem.
func Initialize() bool {
	return C.enet_initialize() == 0
}

// Deinitialize deinitializes the enet subsysetm. This should always be called
// when enet is no longer needed, or sockets and other resources could be leaked.
func Deinitialize() {
	C.enet_deinitialize()
}
