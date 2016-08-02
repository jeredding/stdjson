package stdjson

import (
	"github.com/nkvoll/stdjson/testutil"
	. "github.com/smartystreets/goconvey/convey"
	"os/exec"
	"syscall"
	"testing"
)

func TestProcess(t *testing.T) {
	ct := testutil.NewCapturingTest(t, true, nil)

	Convey("should run echo successfully", ct, func() {
		p := NewProcess("echo", "hello, world")

		err := p.Run()
		So(err, ShouldBeNil)
	})

	Convey("should forward signals successfully", ct, func() {
		p := NewProcess("sleep", "10")

		runErr := make(chan error)
		go func() {
			err := p.Run()
			runErr <- err
		}()

		p.WaitUntilRunning()

		// send the SIGTERM signal to the sleep command, which should result in it exiting.
		err := p.signal(syscall.SIGTERM)
		So(err, ShouldBeNil)

		err = <-runErr
		So(err, ShouldNotBeNil)

		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				So(status.ExitStatus(), ShouldEqual, -1)
			} else {
				t.Errorf("unknown exit code: %s", exiterr)
			}
		}
	})
}
