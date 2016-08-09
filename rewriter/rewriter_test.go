package rewriter

import (
	"github.com/nkvoll/stdjson/config"
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"testing"
)

func TestRewriterForStreamConfig(t *testing.T) {
	Convey("should be able to create rewriters for example configurations", t, func() {
		exampleFiles := []string{
			"../examples/default-fields.yaml",
			"../examples/ls-rewriter.yaml",
			"../examples/ls-advanced.yaml",
			"../examples/noop.yaml",
		}

		for _, exampleFile := range exampleFiles {
			c, err := config.LoadConfigFromFile(exampleFile)
			So(err, ShouldBeNil)
			So(c, ShouldNotBeNil)

			if c.Stdout != nil {
				r, err := RewriterForStreamConfig(c.Stdout, os.Stdout)
				So(err, ShouldBeNil)
				So(r, ShouldNotBeNil)
			}

			if c.Stderr != nil {
				r, err := RewriterForStreamConfig(c.Stderr, os.Stderr)
				So(err, ShouldBeNil)
				So(r, ShouldNotBeNil)
			}
		}
	})
}
