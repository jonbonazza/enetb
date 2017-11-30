#include "enet/enet.h"
#include "enet/enetb.h"

ENetAddress* create_enet_address(const char* host, enet_uint16 port) {
    ENetAddress* addr = enet_malloc(sizeof(ENetAddress));
    int ret = enet_address_set_host(addr, host);
    if (ret > 0) {
        enet_free(addr);
        return NULL;
    }
    addr->port = port;
    return addr;
}