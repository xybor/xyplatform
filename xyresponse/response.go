package xyresponse

import (
	"fmt"
)

type metadata struct {
	service string
	version float64
	errno int
	errer string
}

type xyresponse struct {
	data string
	meta metadata
}

func NewGenerator(serv string, ver float64) xyresponse {
	return xyresponse{
		data: "",
		meta: {
			service: serv,
			version: ver,
			errno: 0,
		},
	}
}

func (x *xyresponse) New(new_dt string) xyresponse {
	x.data = new_dt
}

func (x *xyresponse) NewError(new_err string) xyresponse {
	x.errer = new_err
}

func (x *xyresponse) New(new_dt string, new_mt metadata){
	x.data = new_dt,
	x.meta = new_mt
}

