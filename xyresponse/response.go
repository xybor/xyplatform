package xyresponse

type xyresponse struct {
	data any
	meta map[string]any
}

func NewGenerator(serv string, ver float64) xyresponse {
	return xyresponse{
		data: "",
		meta: map[string]any{
			"service": serv,
			"version": ver,
			"errno":   0,
			"error":   "",
		},
	}
}

func (x *xyresponse) NewData(new_dt any) xyresponse {
	x.data = new_dt
	return *x
}

func (x *xyresponse) NewError(new_err error) xyresponse {
	x.meta["error"] = new_err
	return *x
}

func (x *xyresponse) NewMeta(new_key string, new_value any) {
	x.meta[new_key] = new_value
	return
}
