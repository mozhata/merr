package merr

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

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
	convey.Convey("test wrapErr", t, func() {
		convey.Convey("test err is nil", func() {
			msgs := "err is nil"
			e := wrapErr(nil, originCode, msgs)
			convey.So(e.Error(), convey.ShouldEqual, msgs)
			convey.So(e.Code, convey.ShouldEqual, originCode)
			convey.So(strings.Contains(e.CallStack(), "github.com/mozhata/merr/error_test.go"), convey.ShouldBeTrue)
		})
		convey.Convey("err not nil, masgs is emtpy", func() {
			e := wrapErr(originErr, originCode)
			convey.So(e.Error(), convey.ShouldEqual, originErr.Error())
			convey.So(e.Code, convey.ShouldEqual, originCode)
		})
		convey.Convey("err and masgs not emtpy", func() {
			msg := "msg"
			e := wrapErr(originErr, originCode, msg)
			convey.So(e.Error(), convey.ShouldEqual, msg)
			convey.So(e.RawErr(), convey.ShouldEqual, originErr)
			convey.So(e.Code, convey.ShouldEqual, originCode)
		})
		convey.Convey("wrap Merr", func() {
			msgv1 := "msg v1"
			msgv2 := "msg v2"
			codev2 := 2
			codev3 := 0

			err := wrapErr(originErr, originCode, msgv1)
			convey.So(err.Error(), convey.ShouldEqual, msgv1)
			convey.So(err.RawErr(), convey.ShouldEqual, originErr)
			convey.So(err.Code, convey.ShouldEqual, originCode)

			err = wrapErr(err, codev2, msgv2)
			convey.So(err.Error(), convey.ShouldEqual, msgv2)
			convey.So(err.RawErr(), convey.ShouldEqual, originErr)
			convey.So(err.Code, convey.ShouldEqual, codev2)

			err = wrapErr(err, codev3)
			convey.So(err.Code, convey.ShouldEqual, codev2)
			convey.So(err.Error(), convey.ShouldEqual, msgv2)
		})
	})
	convey.Convey("test WrapErr", t, func() {
		e := WrapErr(originErr)
		convey.So(e.RawErr(), convey.ShouldEqual, originErr)
		convey.So(e.Code, convey.ShouldEqual, http.StatusInternalServerError)
	})
	convey.Convey("test WrapErrWithCode", t, func() {
		e := WrapErrWithCode(originErr, originCode)
		convey.So(e.RawErr(), convey.ShouldEqual, originErr)
		convey.So(e.Code, convey.ShouldEqual, originCode)
	})
	convey.Convey("test InternalErr", t, func() {
		err := InternalErr(originErr)
		e, ok := err.(*MErr)
		convey.So(ok, convey.ShouldBeTrue)
		convey.So(e.RawErr(), convey.ShouldEqual, originErr)
		convey.So(e.Code, convey.ShouldEqual, http.StatusInternalServerError)
	})
	convey.Convey("test ForbiddenErr", t, func() {
		err := ForbiddenErr(originErr)
		e, ok := err.(*MErr)
		convey.So(ok, convey.ShouldBeTrue)
		convey.So(e.RawErr(), convey.ShouldEqual, originErr)
		convey.So(e.Code, convey.ShouldEqual, http.StatusForbidden)
	})
	convey.Convey("test InvalidErr", t, func() {
		err := InvalidErr(originErr)
		e, ok := err.(*MErr)
		convey.So(ok, convey.ShouldBeTrue)
		convey.So(e.RawErr(), convey.ShouldEqual, originErr)
		convey.So(e.Code, convey.ShouldEqual, http.StatusBadRequest)
	})
	convey.Convey("test NotFoundErr", t, func() {
		err := NotFoundErr(originErr)
		e, ok := err.(*MErr)
		convey.So(ok, convey.ShouldBeTrue)
		convey.So(e.RawErr(), convey.ShouldEqual, originErr)
		convey.So(e.Code, convey.ShouldEqual, http.StatusNotFound)
	})
}
