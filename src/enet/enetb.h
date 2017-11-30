#ifndef __ENET_ENETB_H__
#define __ENET_ENETB_H__
#ifdef __cplusplus
extern "C"
{
#endif
#include "enet/enet.h"

ENetAddress* create_enet_address(const char* host, enet_uint16 port);
#ifdef __cplusplus
}
#endif
#endif