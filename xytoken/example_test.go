package xytoken_test

import (
	"fmt"
	"time"

	"github.com/xybor/xyplatform/xyerror"
	"github.com/xybor/xyplatform/xytoken"
)

func ExampleNewEngine() {
	type User struct {
		Name string
		Role string
	}

	jwtDriver := xytoken.NewJWTSecretDriver("SECRET")
	engine := xytoken.NewEngine(jwtDriver, time.Hour)

	user := User{"John Doe", "admin"}
	validToken, err := engine.Generate(user)
	if err != xyerror.Success {
		fmt.Println(err)
	}

	var parseUser User
	err = engine.Parse(validToken, &parseUser)
	if err != xyerror.Success {
		fmt.Println(err)
	}
	fmt.Println(parseUser)

	engine.SetExpiration(-2 * time.Hour)
	expriredToken, err := engine.Generate(user)
	if err != xyerror.Success {
		fmt.Println(err)
	}

	err = engine.Parse(expriredToken, &parseUser)
	if err != xyerror.Success {
		fmt.Println(err)
	}

	// Output:
	// {John Doe admin}
	// [60008][ExpiredError] Token expired
}
