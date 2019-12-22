package merr

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

// MErr basic error class
type MErr struct {
	Msg     string    // 对应Error()
	Code    int       // 错误码
	rawErr  error     // 初始错误信息, 不会被更新
	stackPC []uintptr // 初始错误的调用栈
}

func (e *MErr) Error() string {
	return e.Msg
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
		// TODO: make configable
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

// ErrDetail get detail error info
func ErrDetail(err error) string {
	e := Wrap(err, 0)
	return fmt.Sprintf("E%d: err: %s\nraw err: %s\ncall stack: %s\n",
		e.Code,
		e.Error(),
		e.RawErr(),
		e.CallStack(),
	)
}

// WrapDefaultCode equal to UnKnowErr(err)
// use Wrap func return *MErr value.
// notice: be careful, the returned value is *MErr, not error.
func WrapDefaultCode(err error, fmtAndArgs ...interface{}) *MErr {
	return Wrap(err, 0, fmtAndArgs...)
}

// NilOrWrap return nil if err param is nil, otherwise, then Wrap
// notice: be careful, the returned value is *MErr, not error.
func NilOrWrap(err error, code int, fmtAndArgs ...interface{}) error {
	if err == nil {
		return nil
	}
	return Wrap(err, code, fmtAndArgs...)
}

// Wrap notice: be careful, the returned value is *MErr, not error
func Wrap(err error, code int, fmtAndArgs ...interface{}) *MErr {
	return WrapDepth(1, err, code, fmtAndArgs...)
}

// WrapDepth The argument depth is the number of stack frames to skip before recording in pc,
// with 0 identifying the caller of WrapDepth.
// if a wrapper is added to WrapDepth, depth should +1, like Wrap
func WrapDepth(depth int, err error, code int, fmtAndArgs ...interface{}) *MErr {
	msg := fmtErrMsg(fmtAndArgs...)
	if err == nil {
		err = errors.New(msg)
	}
	if e, ok := err.(*MErr); ok {
		if msg != "" {
			e.Msg = msg
		}
		if code != 0 {
			e.Code = code
		}
		return e
	}

	pcs := make([]uintptr, 32)
	// skip some first invocations
	count := runtime.Callers(2+depth, pcs)
	e := &MErr{
		Code:    code,
		Msg:     msg,
		rawErr:  err,
		stackPC: pcs[:count],
	}
	if e.Msg == "" {
		e.Msg = err.Error()
	}
	return e
}

// fmtErrMsg used to format error message
func fmtErrMsg(msgs ...interface{}) string {
	if len(msgs) > 1 {
		return fmt.Sprintf(msgs[0].(string), msgs[1:]...)
	}
	if len(msgs) == 1 {
		if v, ok := msgs[0].(string); ok {
			return v
		}
		if v, ok := msgs[0].(error); ok {
			return v.Error()
		}
	}
	return ""
}
