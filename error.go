package merr

import (
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"strings"
)

// MErr basic error class
type MErr struct {
	Message    string
	StatusCode int
	rawErr     error
	stackPC    []uintptr
}

// RawErr the origin err
func (e MErr) RawErr() error {
	return e.rawErr
}

// CallStack get function call stack
func (e MErr) CallStack() string {
	frames := runtime.CallersFrames(e.stackPC)
	var (
		f      runtime.Frame
		more   bool
		result string
		index  int
	)
	for {
		f, more = frames.Next()
		if index = strings.Index(f.File, "src"); index != -1 {
			// trim GOPATH or GOROOT prifix
			f.File = string(f.File[index+4:])
		}
		result = fmt.Sprintf("%s%s\n\t%s:%d\n", result, f.Function, f.File, f.Line)
		if !more {
			break
		}
	}
	return result
}

func (e *MErr) Error() string {
	return e.Message
}

// NotFoundErr use http.StatusNotFound to express not found err
func NotFoundErr(err error, fmtAndArgs ...interface{}) error {
	return wrapErr(err, http.StatusNotFound, fmtAndArgs...)
}

// InvalidErr use http.StatusBadRequest to express bad params err
func InvalidErr(err error, fmtAndArgs ...interface{}) error {
	return wrapErr(err, http.StatusBadRequest, fmtAndArgs...)
}

// ForbiddenErr use http.StatusForbidden to express permission deny err
func ForbiddenErr(err error, fmtAndArgs ...interface{}) error {
	return wrapErr(err, http.StatusForbidden, fmtAndArgs...)
}

// InternalErr use http.StatusInternalServerError to express internal server err
func InternalErr(err error, fmtAndArgs ...interface{}) error {
	return wrapErr(err, http.StatusInternalServerError, fmtAndArgs...)
}

// WrapErr equal to InternalErr(err)
func WrapErr(err error, fmtAndArgs ...interface{}) *MErr {
	return wrapErr(err, http.StatusInternalServerError, fmtAndArgs...)
}

func WrapErrWithCode(err error, code int, fmtAndArgs ...interface{}) *MErr {
	return wrapErr(err, code, fmtAndArgs...)
}

// maintain rawErr and update Message if fmtAndArgs is not empty
// notice: the returned value is used as error, so, should not return nil
func wrapErr(err error, code int, fmtAndArgs ...interface{}) *MErr {
	msg := BuildErrMsg(fmtAndArgs...)
	if err == nil {
		err = errors.New(msg)
	}
	if e, ok := err.(*MErr); ok {
		if msg != "" {
			e.Message = msg
		}
		return e
	}

	pcs := make([]uintptr, 32)
	// skip the first 3 invocations
	count := runtime.Callers(3, pcs)
	e := &MErr{
		StatusCode: code,
		Message:    msg,
		rawErr:     err,
		stackPC:    pcs[:count],
	}
	if e.Message == "" {
		e.Message = err.Error()
	}
	return e
}
