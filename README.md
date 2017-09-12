# merr
wrap golang error to make it easy to debug

# usage

```
package main

import (
	"errors"
	"fmt"

	"github.com/mozhata/merr"
)

func main() {
	var err error
	err = merr.WrapErr(nil, "this is an error")
	logMerr(err)
}

func warpperD() error {
	err := errors.New("origin err")
	return merr.WrapErr(err)
}

func wrapperA(err error) error {
	return merr.WrapErr(err, "wrap err by wrapper A")
}

func warpperB(err error) error {
	return merr.WrapErrWithCode(err, 123, "wrap err by wrapper B")
}

func warpperC(err error) error {
	return merr.WrapErr(err)
}

func logMerr(err error) {
	e := merr.WrapErr(err)
	fmt.Printf("E%d: err: %s\nraw err: %s\ncall stack: %s\n",
		e.StatusCode,
		e.Error(),
		e.RawErr(),
		e.CallStack(),
	)
}

```
output:
```
E500: err: this is an error
raw err: this is an error
call stack: main.main
	practice/go/example/main.go:12
runtime.main
	runtime/proc.go:183
runtime.goexit
	runtime/asm_amd64.s:2086

```
func main() {
	var err error
	err = merr.WrapErr(nil, "this is an error")
	err = wrapperA(err)
	logMerr(err)
}


```

```
output:
```
E500: err: wrap err by wrapper A
raw err: this is an error
call stack: main.main
	practice/go/example/main.go:12
runtime.main
	runtime/proc.go:183
runtime.goexit
	runtime/asm_amd64.s:2086

```

```
func main() {
	var err error
	err = merr.WrapErr(nil, "this is an error")
	err = wrapperA(err)
	err = warpperB(err)
	err = warpperC(err)
	logMerr(err)
}

```
output: 
```
E500: err: wrap err by wrapper B
raw err: this is an error
call stack: main.main
	practice/go/example/main.go:12
runtime.main
	runtime/proc.go:183
runtime.goexit
	runtime/asm_amd64.s:2086

```

