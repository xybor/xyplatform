package xypacket

type Packet struct {
	Data *interface{}            `json:"data,omitempty"`
	Meta *map[string]interface{} `json:"meta,omitempty"`
}

func CreatePacket(
	data *interface{},
	meta *map[string]interface{},
) Packet {
	return Packet{Data: data, Meta: meta}
}
