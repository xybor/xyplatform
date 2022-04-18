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
	for _, c := range chans {
		eselector.Recv(xyselect.C(c))
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
	for _, c := range chans {
		rselector.Recv(xyselect.C(c))
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
