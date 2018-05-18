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
	err = warpper()
	logMerr(err)
}

func warpper() error {
	err := errors.New("origin err")
	return merr.WrapErr(err, "new err")
}

func logMerr(err error) {
	e := merr.WrapErr(err)
	fmt.Printf("E%d: err: %s\nraw err: %s\ncall stack: %s\n",
		e.Code,
		e.Error(),
		e.RawErr(),
		e.CallStack(),
	)
}

```

output:

```
E500: err: new err
raw err: origin err
call stack: main.warpper
	practice/go/example/main.go:18
main.main
	practice/go/example/main.go:12
runtime.main
	runtime/proc.go:183
runtime.goexit
	runtime/asm_amd64.s:2086

```

