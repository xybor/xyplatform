package xytoken

import (
	"errors"
	"time"
)

type Payload struct {
	IssuedAt  time.Time
	ExpiresAt time.Time
	Data      interface{}
}

func (p *Payload) Valid() error {
	if time.Now().After(p.ExpiresAt) {
		return errors.New("token expired")
	}
	return nil
}

func newPayload(expiration time.Duration, data interface{}) Payload {
	now := time.Now()
	return Payload{
		IssuedAt:  now,
		ExpiresAt: now.Add(expiration),
		Data:      data,
	}
}
