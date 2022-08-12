package xyselect_test

import (
	"errors"
	"fmt"

	"github.com/xybor/xyplatform/xyselect"
)

func ExampleE() {
	var chans []chan int
	for i := 0; i < 10; i++ {
		chans = append(chans, make(chan int))
		go func(c chan int, v int) {
			c <- v
			close(c)
		}(chans[i], i)
	}

	eselector := xyselect.E()
	for i := range chans {
		eselector.Recv(xyselect.C(chans[i]))
	}

	okCounter := 0
	for {
		_, _, err := eselector.Select(false)

		if !errors.Is(err, xyselect.ExhaustedError) {
			okCounter += 1
		} else {
			break
		}
	}
	fmt.Println(okCounter)

	// Output:
	// 10
}

func ExampleR() {
	var chans []chan int
	for i := 0; i < 10; i++ {
		chans = append(chans, make(chan int))
		go func(c chan<- int, v int) {
			c <- v
			close(c)
		}(chans[i], i)
	}

	rselector := xyselect.R()
	for i := range chans {
		rselector.Recv(xyselect.C(chans[i]))
		rselector.Send(make(chan int, 1), 1)
	}

	okCounter := 0
	for {
		_, _, err := rselector.Select(false)

		if err != nil {
			okCounter += 1
		}

		if okCounter == 20 {
			break
		}
	}

	// Output:
}
