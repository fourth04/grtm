# about grtm
[![Build Status](https://travis-ci.org/scottkiss/grtm.svg?branch=master)](https://travis-ci.org/scottkiss/grtm)

grtm is a tool to manage golang goroutines.use this can start or stop a long loop goroutine.

## Getting started
```bash
go get github.com/fourth04/grtm
```

## Create normal goroutine

```go
package main

import (
        "fmt"
        "github.com/fourth04/grtm"
        "time"
       )

func normal() {
    fmt.Println("i am normal goroutine")
    time.Sleep(time.Second * time.Duration(5))
}

func main() {
        gm := &grtm.GrManager{}
        gm.NewNormalGoroutine("normal", normal)
        fmt.Println("main function")
        gm.Wait()
}
```

## Create normal goroutine function with params

```go
package main

import (
	"fmt"
	"time"

	"github.com/fourth04/grtm"
)

func normal() {
	fmt.Println("i am normal goroutine")
}

func funcWithParams(args ...interface{}) {
	fmt.Println(args[0].([]interface{})[0].(string))
	fmt.Println(args[0].([]interface{})[1].(string))
	time.Sleep(time.Second * time.Duration(5))
}

func main() {
	gm := &grtm.GrManager{}
	gm.NewNormalGoroutine("normal", normal)
	fmt.Println("main function")
	gm.NewNormalGoroutine("funcWithParams", funcWithParams, "hello", "world")
    gm.Wait()
}
```

## Create long loop goroutine then stop it

```go
package main

import (
        "fmt"
        "github.com/fourth04/grtm"
        "time"
       )

func myfunc() {
	fmt.Println("do something repeat by interval 4 seconds")
	time.Sleep(time.Second * time.Duration(4))
}

func main() {
	gm := &grtm.GrManager{}
	gm.NewLoopGoroutine("myfunc", myfunc)
	fmt.Println("main function")
	time.Sleep(time.Second * time.Duration(20))
	fmt.Println("stop myfunc goroutine")
	gm.StopGoroutine("myfunc")
    gm.Wait()
}
```

output

```bash
main function
no signal
do something repeat by interval 4 seconds
no signal
do something repeat by interval 4 seconds
no signal
do something repeat by interval 4 seconds
no signal
do something repeat by interval 4 seconds
no signal
do something repeat by interval 4 seconds
stop myfunc goroutine
[gid: 5577006791947779410, name: myfunc] quit
```

## Create diy goroutine

```go
package main

import (
	"fmt"
	"time"

	"github.com/fourth04/grtm"
)

func diy(chMsg chan string) {
	for {
		select {
		case <-chMsg:
			return
		default:
			fmt.Println("no signal")
		}
		fmt.Println("i am diy goroutine")
		time.Sleep(time.Second * time.Duration(1))
	}
}

func main() {
	gm := &grtm.GrManager{}
	gm.NewDiyGoroutine("diy", diy)
	fmt.Println("main function")
	time.Sleep(time.Second * time.Duration(5))
    gm.StopDiyGoroutine("diy")
    gm.Wait()
}
```

## Create diy goroutine function with params

```go
package main

import (
	"fmt"
	"time"

	"github.com/fourth04/grtm"
)

func funcWithParams(chMsg chan string, args ...interface{}) {
	for {
		select {
		case <-chMsg:
			return
		default:
            fmt.Println("no signal")
            fmt.Println(args[0].([]interface{})[0].(string))
            fmt.Println(args[0].([]interface{})[1].(string))
            time.Sleep(time.Second * time.Duration(1))
        }
	}
}

func main() {
	gm := &grtm.GrManager{}
	gm.NewDiyGoroutine("funcWithParams", funcWithParams, "hello", "world")
	fmt.Println("main function")
	time.Sleep(time.Second * time.Duration(5))
	gm.StopDiyGoroutine("funcWithParams")
    gm.Wait()
}
```
