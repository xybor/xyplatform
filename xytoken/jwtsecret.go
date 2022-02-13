package xytoken

import (
	"errors"

	"github.com/golang-jwt/jwt"
	"github.com/mitchellh/mapstructure"
	"github.com/xybor/xyplatform/xyerror"
)

type JWTSecretDriver struct {
	secret string
}

func NewJWTSecretDriver(secret string) JWTSecretDriver {
	return JWTSecretDriver{
		secret: secret,
	}
}

func (d *JWTSecretDriver) SetSecret(secret string) {
	d.secret = secret
}

func (d JWTSecretDriver) generate(payload Payload) (string, xyerror.XyError) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &payload)
	tokenString, err := token.SignedString([]byte(d.secret))

	if err != nil {
		return "", CanNotCreateError.New("Can not create token")
	} else {
		return tokenString, xyerror.Success
	}
}

func (d JWTSecretDriver) parse(token string, output interface{}) xyerror.XyError {
	parser, err := jwt.ParseWithClaims(token, &Payload{}, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(d.secret), nil
	})

	if err != nil || !parser.Valid {
		if err.Error() == jwt.ErrSignatureInvalid.Error() {
			return SignatureInvalidError.New("Invalid signature")
		} else if err.Error() == "token expired" {
			return ExpiredError.New("Token expired")
		} else {
			return InvalidError.New("Invalid token")
		}
	}

	payload, ok := parser.Claims.(*Payload)
	if !ok {
		return ParseError.New("Can not parse token")
	}

	err = mapstructure.Decode(payload.Data, &output)
	if err != nil {
		return ParseError.New("Can not parse token")
	}

	return xyerror.Success
}
