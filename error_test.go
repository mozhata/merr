package merr

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestShortenPath(t *testing.T) {
	convey.Convey("test shortenPath", t, func() {
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
			convey.So(output, convey.ShouldEqual, u.expect)
		}
	})
}
func TestFmtErrMsg(t *testing.T) {
	convey.Convey("test fmtErrMsg", t, func() {
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
			convey.So(output, convey.ShouldEqual, u.expect)
		}
	})
}

func TestMErr(t *testing.T) {
	originErr := errors.New("origin error")
	originCode := 1
	convey.Convey("test Wrap", t, func() {
		convey.Convey("test err is nil", func() {
			msgs := "err is nil"
			e := Wrap(nil, originCode, msgs)
			convey.So(e.Error(), convey.ShouldEqual, msgs)
			convey.So(e.Code, convey.ShouldEqual, originCode)
			convey.So(strings.Contains(e.CallStack(), "merr/error_test.go"), convey.ShouldBeTrue)
		})
		convey.Convey("err not nil, masgs is emtpy", func() {
			e := Wrap(originErr, originCode)
			convey.So(e.Error(), convey.ShouldEqual, originErr.Error())
			convey.So(e.Code, convey.ShouldEqual, originCode)
		})
		convey.Convey("err and masgs not emtpy", func() {
			msg := "msg"
			e := Wrap(originErr, originCode, msg)
			convey.So(e.Error(), convey.ShouldEqual, msg)
			convey.So(e.RawErr(), convey.ShouldEqual, originErr)
			convey.So(e.Code, convey.ShouldEqual, originCode)
		})
		convey.Convey("err, code, msg all empty", func() {
			var err error = Wrap(nil, 0)
			e := Wrap(nil, 0)

			convey.So(e != nil, convey.ShouldBeTrue)
			convey.So(err != nil, convey.ShouldBeTrue)
			err = NilOrWrap(nil, 0)
			convey.So(err == nil, convey.ShouldBeTrue)
		})
		convey.Convey("wrap Merr", func() {
			msgv1 := "msg v1"
			msgv2 := "msg v2"
			codev2 := 2
			codev3 := 0

			err := Wrap(originErr, originCode, msgv1)
			convey.So(err.Error(), convey.ShouldEqual, msgv1)
			convey.So(err.RawErr(), convey.ShouldEqual, originErr)
			convey.So(err.Code, convey.ShouldEqual, originCode)

			err = Wrap(err, codev2, msgv2)
			convey.So(err.Error(), convey.ShouldEqual, msgv2)
			convey.So(err.RawErr(), convey.ShouldEqual, originErr)
			convey.So(err.Code, convey.ShouldEqual, codev2)
			// code, msg both empty
			err = Wrap(err, codev3)
			convey.So(err.Error(), convey.ShouldEqual, msgv2)
			convey.So(err.RawErr(), convey.ShouldEqual, originErr)
			convey.So(err.Code, convey.ShouldEqual, codev2)
			// update code, msg empty
			err = Wrap(err, 22)
			convey.So(err.Error(), convey.ShouldEqual, msgv2)
			convey.So(err.RawErr(), convey.ShouldEqual, originErr)
			convey.So(err.Code, convey.ShouldEqual, 22)
			// update msg empty, code emtpy
			err = Wrap(err, 0, "new err msg")
			convey.So(err.Error(), convey.ShouldEqual, "new err msg")
			convey.So(err.RawErr(), convey.ShouldEqual, originErr)
			convey.So(err.Code, convey.ShouldEqual, 22)
		})

	})
}

func TestMerrWrap(t *testing.T) {
	// wrapHappendLine mast alongside with the first Wrap
	originErr := errors.New("this is origin err")
	var wrapHappendLine int
	err1 := Wrap(originErr, 1, "err1")
	err2 := Wrap(err1, 2, "err2")
	err3 := Wrap(err2, 3, "err3")
	_, _, l, ok := runtime.Caller(3)
	if ok {
		wrapHappendLine = l + 1
		submsg := fmt.Sprintf("merr/error_test.go:%d", wrapHappendLine)
		fmt.Printf("wrapHappendLine: %d\n", wrapHappendLine)
		callstack := err3.CallStack()
		if !strings.Contains(callstack, submsg) {
			t.Fatalf("callstack %s expected contains %s", callstack, submsg)
		}
	}

}
