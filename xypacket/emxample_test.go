package xypacket_test

import (
	"github.com/xybor/xyplatform/xyerror"
	"github.com/xybor/xyplatform/xypacket"
)

// Summary:
//			Packet is used for communicating between client and server.
//	A Packet store data and meta. Data contains main data, meta contains
//	information such as error message (errmsg), error number (errno)
//	and API version.

// Create Example Success packet
func ExampleSuccessPacket(data interface{}) xypacket.Packet {
	meta := map[string]interface{}{
		"APIVersion": 1,
		"errno":      xyerror.Success.Errno(),
		"errmsg":     xyerror.Success.Errmsg(),
	}

	return xypacket.CreatePacket(&data, &meta)
}

// Create Packet with IOError message
func ExampleIOErrorPacket(data interface{}) xypacket.Packet {
	meta := map[string]interface{}{
		"APIVersion": 1,
		"errno":      xyerror.IOError.Errno(),
		"errmsg":     xyerror.IOError.Errno(),
	}

	return xypacket.CreatePacket(&data, &meta)
}
