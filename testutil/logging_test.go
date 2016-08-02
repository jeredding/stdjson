package testutil

import (
	"bytes"
	"github.com/Sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestCapturingTest(t *testing.T) {
	Convey("should be able to capture output from logrus", t, func() {
		t2 := &testing.T{}
		ct := NewCapturingTest(t2, true, nil)

		logrus.Infoln("hello from logrus")

		So(ct.Writer.buffer.String(), ShouldContainSubstring, "hello from logrus")
	})

	Convey("should not dump captured output when the test does not fail", t, func() {
		out := &bytes.Buffer{}
		t2 := &testing.T{}
		_ = NewCapturingTest(t2, true, out)

		logrus.Infoln("hello from logrus")

		So(out.String(), ShouldBeEmpty)
	})

	Convey("should dump captured output when the test fails", t, func() {
		out := &bytes.Buffer{}
		t2 := &testing.T{}
		ct := NewCapturingTest(t2, true, out)

		logrus.Infoln("hello from logrus")

		ct.Fail()
		So(out.String(), ShouldContainSubstring, "hello from logrus")
	})
}
