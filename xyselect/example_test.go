package xyselect_test

import (
	"fmt"

	"github.com/xybor/xyplatform/xyselect"
)

func ExampleE() {
	var c = make(chan int)
	go func() {
		c <- 10
		close(c)
	}()

	var eselector = xyselect.E()
	eselector.Recv(xyselect.C(c))

	var _, v, _ = eselector.Select(false)
	fmt.Println(v)

	// Output:
	// 10
}

func ExampleR() {
	var rselector = xyselect.R()
	var c = make(chan int)
	var rc = xyselect.C(c)

	go func() { c <- 10 }()
	rselector.Recv(rc)
	var _, v, _ = rselector.Select(false)
	fmt.Println("receive", v)

	rselector.Send(c, 20)
	rselector.Select(false)
	fmt.Println("send", <-rc)

	// Output:
	// receive 10
	// send 20
}
