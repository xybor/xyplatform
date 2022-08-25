# Introduction

Xyselect is a library used to call `select` with an unknown number of `case`
statements.

# Features

The main object in the library is `Selector`, a custom usage of `select`
statement.

There are two types of `Selector`, which are `R` (stand for `reflect`) and `E`
(stand for exhausted).

`E` selector uses a center channel to receive all selected channels, so that it
doesn't support to wait on a sending channel. Moreover, `E` selector creates a
goroutine to wait on the center channel, and that goroutine only stops when all
selected channels closed. For this reason, you should call `Select` until get
an error of `ExhaustedError` to ensure the goroutine stopped.

`R` selector uses the built-in library, `reflect`, to customize `select`
statement. It supports to wait on both receiving and sending channels. It also
does not create any goroutine while using.

Each selector has its own advantage, while `R` selector more flexible, `E`
selector is faster.

Visit [pkg.go.dev](https://pkg.go.dev/github.com/xybor/xyplatform/xyselect) for
more details.

# Benchmark

| op name   | time per op |
| --------- | ----------- |
| RSelector | 728ns       |
| ESelector | 679ns       |

# Example

```golang
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
```

```golang
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
```
