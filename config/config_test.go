package config

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestLoadExamplesConfig(t *testing.T) {
	Convey("should be able to load default example configurations", t, func() {
		exampleFiles := []string{
			"../examples/noop.yaml",
			"../examples/ls-rewriter.yaml",
			"../examples/ls-advanced.yaml",
		}

		for _, exampleFile := range exampleFiles {
			c, err := LoadConfigFromFile(exampleFile)
			So(err, ShouldBeNil)
			So(c, ShouldNotBeNil)
		}
	})
}

func TestMultilineConfig_MustTimeoutDuration(t *testing.T) {
	Convey("when parsing timeouts", t, func() {
		c := &MultilineConfig{}

		Convey("should parse a valid timeout", func() {
			timeout := "100ms"
			c.Timeout = &timeout

			expected := 100 * time.Millisecond
			So(c.MustTimeoutDuration(), ShouldResemble, &expected)
		})

		Convey("should panic on invalid timeout", func() {
			timeout := "100meters"
			c.Timeout = &timeout

			So(func() { c.MustTimeoutDuration() }, ShouldPanic)
		})

		Convey("should pass through nil timeouts", func() {
			So(c.MustTimeoutDuration(), ShouldBeNil)
		})
	})
}
