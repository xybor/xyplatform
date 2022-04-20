package xyresponse_test

import (
	"github.com/xybor/xyplatform/xyresponse"
)

var ExampleResponse = xyresponse.NewGenerator("XyAuth", 1)

func ExampleNewData() {
	new_data := "new data"
	ExampleResponse.NewData(new_data)
	return
}
