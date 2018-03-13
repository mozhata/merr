package merr

import "fmt"

func LogErr(err error) {
	e := WrapErr(err)
	fmt.Printf("E%d: err: %s\nraw err: %s\ncall stack: %s\n",
		e.StatusCode,
		e.Error(),
		e.RawErr(),
		e.CallStack(),
	)
}
