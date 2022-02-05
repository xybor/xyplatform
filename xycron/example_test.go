package xycron_test

import (
	"fmt"
	"time"

	"github.com/xybor/xyplatform"
	"github.com/xybor/xyplatform/xycron"
	"github.com/xybor/xyplatform/xylog"
)

var _ = xylog.Config(xyplatform.XyCron, xylog.NoAllow())

func Example() {
	scheduler := xycron.New()
	go scheduler.Start()

	scheduler.EverySecond().Once().Params("Hello, world").Do(fmt.Println)
	scheduler.After(time.Second).Params("Hello, world").Do(fmt.Println)

	time.Sleep(2 * time.Second)
	var err = scheduler.Stop()
	if err != nil {
		fmt.Println(err)
	}

	// Output:
	// Hello, world
	// Hello, world
}
func ExampleGlobal() {
	go xycron.Global().Start()

	xycron.Global().EverySecond().Once().Params("Hello, world").Do(fmt.Println)
	xycron.Global().After(time.Second).Params("Hello, world").Do(fmt.Println)

	time.Sleep(2 * time.Second)
	var err = xycron.Global().Stop()
	if err != nil {
		fmt.Println(err)
	}

	// Output:
	// Hello, world
	// Hello, world
}
