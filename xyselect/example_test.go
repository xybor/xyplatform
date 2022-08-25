package xyselect_test

import (
	"errors"
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

	var _, _, err = eselector.Select(false)
	if errors.Is(err, xyselect.ClosedChannelError) {
		fmt.Println(err)
	}

	_, _, err = eselector.Select(false)
	if errors.Is(err, xyselect.ExhaustedError) {
		fmt.Println(err)
	}

	// Output:
	// 10
	// ClosedChannelError: channel closed
	// ExhaustedError: selector is exhausted
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
