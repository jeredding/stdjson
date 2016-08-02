package rewriter

import (
	"github.com/nkvoll/stdjson/config"
	"github.com/nkvoll/stdjson/testutil"
	. "github.com/smartystreets/goconvey/convey"
	"sync"
	"testing"
)

func TestMulilineReaderFocus(t *testing.T) {
	ct := testutil.NewCapturingTest(t, true, nil)

	FocusConvey("should gather multiple lines", ct, func() {
		FocusConvey("with no prefix continuations", func(c C) {
			m := NewMultilineReader(&config.MultilineConfig{StripLastDelimiter: true})

			messages := []string{
				"this is a message",
				"this is another message",
			}

			testMultiline(
				c, m,
				messages,
				messages,
			)
		})

		Convey("with a prefix continuation", func(c C) {
			m := NewMultilineReader(&config.MultilineConfig{StripLastDelimiter: true, PrefixContinuations: []string{" "}})

			testMultiline(
				c, m,
				[]string{"hello,", " world"},
				[]string{"hello, world"},
			)
		})

		Convey("as a single line when timing out", func(c C) {
			m := NewMultilineReader(&config.MultilineConfig{StripLastDelimiter: true, PrefixContinuations: []string{" "}})

			testMultiline(
				c, m,
				[]string{"test"},
				[]string{"test"},
			)
		})

		Convey("with multiple prefix continuations", func(c C) {
			m := NewMultilineReader(&config.MultilineConfig{StripLastDelimiter: true, PrefixContinuations: []string{" ", ","}})

			testMultiline(
				c, m,
				[]string{"hello,", " world", "hello", ", world"},
				[]string{"hello, world", "hello, world"},
			)
		})
	})
}

func testMultiline(c C, m *MultilineBuffer, input []string, expectedOutput []string) {
	running := &sync.WaitGroup{}
	running.Add(1)

	go func() {
		m.Run()
		running.Done()
	}()

	for _, line := range input {
		_, err := m.Write([]byte(line))
		c.So(err, ShouldBeNil)
		_, err = m.Write([]byte{'\n'})
		c.So(err, ShouldBeNil)
	}

	err := m.Close()
	c.So(err, ShouldBeNil)

	actual := []string{}
	for line := range m.Lines {
		actual = append(actual, line)
	}

	running.Wait()

	c.So(actual, ShouldResemble, expectedOutput)
}
