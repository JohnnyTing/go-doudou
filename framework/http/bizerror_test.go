package ddhttp_test

import (
	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
	ddhttp "github.com/unionj-cloud/go-doudou/framework/http"
	"testing"
)

func TestWithStatusCode(t *testing.T) {
	Convey("Create a BizError with 401 http status code", t, func() {
		bizError := ddhttp.NewBizError(errors.New("Unauthorised"), ddhttp.WithStatusCode(401))
		So(bizError, ShouldNotBeZeroValue)
		So(bizError.StatusCode, ShouldEqual, 401)

		Convey("Should output Unauthorised", func() {
			So(bizError.String(), ShouldEqual, "Unauthorised")
		})
	})
}

func TestWithErrCode(t *testing.T) {
	Convey("Create a BizError with 100401 business error code", t, func() {
		bizError := ddhttp.NewBizError(errors.New("Unauthorised"), ddhttp.WithErrCode(100401))
		So(bizError, ShouldNotBeZeroValue)
		So(bizError.ErrCode, ShouldEqual, 100401)

		Convey("Should output 100401 Unauthorised", func() {
			So(bizError.String(), ShouldEqual, "100401 Unauthorised")
		})

		Convey("Should have the same output with String()", func() {
			So(bizError.Error(), ShouldEqual, bizError.String())
		})
	})
}
