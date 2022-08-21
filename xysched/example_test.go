package xysched_test

import (
	"fmt"
	"time"

	"github.com/xybor/xyplatform/xysched"
)

func panicFunc() {
	panic("custom panic function")
}

func ExampleTask() {
	var scheduler = xysched.NewScheduler()

	// Example 1: Task is a simple future used for scheduling to run a function.
	var done = make(chan any)
	var future = xysched.NewTask(func(a ...any) {
		fmt.Println(a...)
		close(done)
	}, "1. foo")
	scheduler.Now() <- future
	<-done

	// Example 2: Callback will be run after the task completed.
	done = make(chan any)
	future = xysched.NewTask(fmt.Println, "2. foo foo")
	future.Callback(func() { close(done) })
	scheduler.Now() <- future
	<-done

	// Example 3: Then adds a callback handling returned values of task after
	// task completed.
	done = make(chan any)
	future = xysched.NewTask(fmt.Sprintf, "3. foo %s", "bar")
	future.Then(func(s string) {
		fmt.Println(s)
		close(done)
	})
	scheduler.Now() <- future
	<-done

	// Example 4: Catch adds a callback handling the panicked error of task if
	// the task panicked.
	// NOTE: if task panics a non-error interface, it will be wrapped into
	//       xysched.CallError.
	done = make(chan any)
	future = xysched.NewTask(panicFunc)
	future.Then(func() {
		fmt.Println("4. Then branch")
		close(done)
	})
	future.Catch(func(e error) {
		fmt.Println("4. Catch branch", e)
		close(done)
	})
	scheduler.Now() <- future
	<-done

	scheduler.Stop()

	// Output:
	// 1. foo
	// 2. foo foo
	// 3. foo bar
	// 4. Catch branch CallError: custom panic function
}

func ExampleGlobal() {
	// Example 1: You can use the global scheduler throughout program without
	// creating a new one.
	var done = make(chan any)
	var future = xysched.NewTask(func() {
		fmt.Println("1. bar bar")
		close(done)
	})
	xysched.Global().Now() <- future
	<-done

	// Example 2: Scheduler can schedule one future After or At a time.
	done = make(chan any)
	future = xysched.NewTask(func() {
		fmt.Println("2. barfoo")
		close(done)
	})
	xysched.Global().After(time.Second) <- future
	<-done

	// Output:
	// 1. bar bar
	// 2. barfoo
}

func wait(c chan any, n int) {
	for i := 0; i < n; i++ {
		<-c
	}
}

func ExampleCron() {
	var scheduler = xysched.NewScheduler()
	// Example 1: Cron is a future which runs function periodically. By default,
	// it runs secondly forever.
	var done = make(chan any)
	var future = xysched.NewCron(func(a ...any) {
		fmt.Println(a...)
		done <- nil
	}, "1.", "foo", "bar")
	scheduler.Now() <- future
	wait(done, 2)
	scheduler.Stop()

	scheduler = xysched.NewScheduler()
	// Example 2: It can modify periodic duration and the maximum times the
	// function could run.
	done = make(chan any)
	future = xysched.NewCron(func() {
		fmt.Println("2. bar bar")
		done <- nil
	}).Every(1 * time.Millisecond).Twice()
	scheduler.Now() <- future
	wait(done, 2)

	// Example 3: Callback, Then, Catch can also be used on cron.
	done = make(chan any)
	future = xysched.NewCron(fmt.Println, "3.", "foobar").Times(3)
	future.Callback(func() { done <- nil })
	scheduler.Now() <- future
	wait(done, 3)

	// Example 4: Finish adds a callback future which will be run when cron ran
	// out of times.
	done = make(chan any)
	future = xysched.NewCron(fmt.Println, "4.", "foobar").Twice()
	future.Finish(func() { close(done) })
	scheduler.Now() <- future
	wait(done, 1)

	scheduler.Stop()

	// Output:
	// 1. foo bar
	// 1. foo bar
	// 2. bar bar
	// 2. bar bar
	// 3. foobar
	// 3. foobar
	// 3. foobar
	// 4. foobar
	// 4. foobar
}
