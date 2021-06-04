package merr

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
	"testing"
)

func TestShortenPath(t *testing.T) {
	type unit struct {
		input  string
		expect string
	}
	table := []unit{
		{},
		{"/p1/p2/p3/file.toml", "p3/file.toml"},
		{"/p3/file.go", "p3/file.go"},
		{"/file.go", "file.go"},
		{"file.go", "file.go"},
	}
	for _, u := range table {
		output := shortenPath(u.input)
		// convey.So(output, convey.ShouldEqual, u.expect)
		AssertEqual(t, output, u.expect)
	}
}
func TestFmtErrMsg(t *testing.T) {
	type unit struct {
		input  []interface{}
		expect string
	}
	table := []unit{
		unit{},
		unit{
			input:  []interface{}{errors.New("new error")},
			expect: "new error",
		},
		unit{
			input:  []interface{}{"single"},
			expect: "single",
		},
		unit{
			input:  []interface{}{"this is a %s", "format"},
			expect: "this is a format",
		},
	}
	for _, u := range table {
		output := fmtErrMsg(u.input...)
		AssertEqual(t, output, u.expect)
	}
}

func TestMerr_ErrNil(t *testing.T) {
	originCode := 1
	msgs := "err is nil"
	e := Wrap(nil, originCode, msgs)
	AssertEqual(t, e.Error(), msgs)
	AssertEqual(t, e.Code, originCode)
	AssertContainSubStr(t, e.CallStack(), "merr/error_test.go")
}
func TestMerr_ErrNilMsgNil(t *testing.T) {
	var err error = Wrap(nil, 0)
	e := Wrap(nil, 0)

	AssertEqual(t, e != nil, true)
	AssertEqual(t, err != nil, true)
	err = NilOrWrap(nil, 0)
	AssertEqual(t, err, nil)
}
func TestMerr_ErrNotNilMsgNil(t *testing.T) {
	originErr := errors.New("origin error")
	originCode := 1
	e := Wrap(originErr, originCode)
	AssertEqual(t, e.Error(), originErr.Error())
	AssertEqual(t, e.Code, originCode)
}
func TestMerr_ErrNotNilMsgNotNil(t *testing.T) {
	originErr := errors.New("origin error")
	originCode := 1
	msg := "msg"
	e := Wrap(originErr, originCode, msg)
	AssertEqual(t, e.Error(), msg)
	AssertEqual(t, e.Code, originCode)
	AssertEqual(t, e.RawErr(), originErr)
}

func TestMerr_DaylyWrok(t *testing.T) {
	originErr := errors.New("origin error")
	originCode := 1
	msgv1 := "msg v1"
	msgv2 := "msg v2"
	codev2 := 2
	codev3 := 0

	err := Wrap(originErr, originCode, msgv1)
	AssertEqual(t, err.Error(), msgv1)
	AssertEqual(t, err.RawErr(), originErr)
	AssertEqual(t, err.Code, originCode)

	err = Wrap(err, codev2, msgv2)
	AssertEqual(t, err.Error(), msgv2)
	AssertEqual(t, err.RawErr(), originErr)
	AssertEqual(t, err.Code, codev2)
	// code, msg both empty
	err = Wrap(err, codev3)
	AssertEqual(t, err.Error(), msgv2)
	AssertEqual(t, err.RawErr(), originErr)
	AssertEqual(t, err.Code, codev2)
	// update code, msg empty
	err = Wrap(err, 22)
	AssertEqual(t, err.Error(), msgv2)
	AssertEqual(t, err.RawErr(), originErr)
	AssertEqual(t, err.Code, 22)
	// update msg empty, code emtpy
	err = Wrap(err, 0, "new err msg")
	AssertEqual(t, err.Error(), "new err msg")
	AssertEqual(t, err.RawErr(), originErr)
	AssertEqual(t, err.Code, 22)
}

func TestMerrWrap(t *testing.T) {
	// wrapHappendLine mast alongside with the first Wrap
	originErr := errors.New("this is origin err")
	var wrapHappendLine int
	err1 := Wrap(originErr, 1, "err1")
	err2 := Wrap(err1, 2, "err2")
	err3 := Wrap(err2, 3, "err3")
	_, _, l, _ := runtime.Caller(0)
	wrapHappendLine = l - 3
	submsg := fmt.Sprintf("merr/error_test.go:%d", wrapHappendLine)
	callstack := err3.CallStack()
	AssertContainSubStr(t, callstack, submsg)
}

func AssertEqual(t *testing.T, actual, expect interface{}) {
	if actual != expect {
		logf("expect %v (actual) = %v (expect)", actual, expect)
		t.FailNow()
	}
}
func AssertContainSubStr(t *testing.T, longStr, subStr string) {
	if !strings.Contains(longStr, subStr) {
		logf("longStr %s expected contains substr %s", longStr, subStr)
		t.FailNow()
	}
}
func logf(format string, a ...interface{}) {
	_, fileName, line, _ := runtime.Caller(2)
	msg := fmt.Sprintf(format, a...)
	fmt.Printf("%s:%d %s\n", fileName, line, msg)
}
