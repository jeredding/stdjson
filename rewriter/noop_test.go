package rewriter

import (
	"bytes"
	"github.com/nkvoll/stdjson/testutil"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNoop(t *testing.T) {
	ct := testutil.NewCapturingTest(t, true, nil)

	Convey("should forward input to its output", ct, func() {
		w := &bytes.Buffer{}
		n := NewNoop(w)

		runErr := make(chan error)
		go func() {
			err := n.Run()
			runErr <- err
		}()

		msg := []byte("this is a message")
		c, err := n.Write(msg)
		So(err, ShouldBeNil)
		So(c, ShouldEqual, len(msg))

		So(w.Bytes(), ShouldResemble, msg)

		err = n.Close()
		So(err, ShouldBeNil)

		err = <-runErr
		So(err, ShouldBeNil)
	})
}
