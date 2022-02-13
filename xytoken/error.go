package xytoken

import (
	"github.com/xybor/xyplatform"
	"github.com/xybor/xyplatform/xyerror"
)

var _ = xyerror.Register(xyplatform.Xytoken)
var TokenError = xyerror.NewType(xyplatform.Xytoken, "TokenError")
var ValidationError = xyerror.NewType(xyplatform.Xytoken, "ValidationError")

var (
	InvalidError              = TokenError.NewType("InvalidError")
	SignatureInvalidError     = TokenError.NewType("SignatureInvalidError")
	InvalidSigningMethodError = TokenError.NewType("InvalidSigningMethodError")
	ParseError                = TokenError.NewType("ParsePayloadError")
	CanNotCreateError         = TokenError.NewType("CanNotCreateError")

	ExpiredError = ValidationError.NewType("ExpiredError")
)
