package stdjson

import (
	"encoding/json"
	"github.com/Sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestLogFormatter(t *testing.T) {
	Convey("should format times as UTC", t, func() {
		f := &LogFormatter{}

		entry := logrus.Entry{
			Time: time.Time{}.Add(24 * time.Hour).In(time.FixedZone("test", 3600)),
		}

		formatted, err := f.Format(&entry)
		So(err, ShouldBeNil)

		v := map[string]interface{}{}

		err = json.Unmarshal(formatted, &v)
		So(err, ShouldBeNil)

		So(v["time"], ShouldEqual, "0001-01-02T00:00:00Z")
	})
}
