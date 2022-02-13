package xytoken

import (
	"time"

	"github.com/xybor/xyplatform/xyerror"
)

type driver interface {
	generate(Payload) (string, xyerror.XyError)
	parse(string, interface{}) xyerror.XyError
}

type engine struct {
	d          driver
	expiration time.Duration
}

func (e *engine) SetExpiration(expiration time.Duration) {
	e.expiration = expiration
}

func NewEngine(driver driver, expiration time.Duration) engine {
	return engine{
		d:          driver,
		expiration: expiration,
	}
}

func (e engine) Generate(data interface{}) (string, xyerror.XyError) {
	payload := newPayload(e.expiration, data)
	return e.d.generate(payload)
}

func (e engine) Parse(token string, output interface{}) xyerror.XyError {
	return e.d.parse(token, output)
}
