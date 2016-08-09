package rewriter

import (
	"bytes"
	"github.com/nkvoll/stdjson/config"
	"github.com/nkvoll/stdjson/testutil"
	_ "github.com/nkvoll/stdjson/testutil"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/vjeantet/grok"
	"sync"
	"testing"
	"io"
)

func TestMultilineGrokRewriter(t *testing.T) {
	ct := testutil.NewCapturingTest(t, true, nil)

	Convey("should rewrite", ct, func() {
		b := &bytes.Buffer{}
		g, err := grok.NewWithConfig(&grok.Config{
			NamedCapturesOnly: true,
		})
		So(err, ShouldBeNil)
		g.AddPattern("LSTIMESTAMP", "%{MONTH} +%{MONTHDAY} +%{HOUR}:%{MINUTE}")
		g.AddPattern("GREEDYALL", `(\n|.)*`)

		conf := &config.GrokRewriterConfig{
			MatchPatterns: []string{"%{GREEDYDATA:any}"},
			Multiline: config.MultilineConfig{
				StripLastDelimiter: true,
			},
		}

		Convey("should capture input matching a pattern", func(c C) {
			m := mustNewMultiGrokRewriter(c, b, g, conf)
			runRewriter(c, m, []string{"hello"})
			So(b.String(), ShouldEqual, `{"any":"hello"}`+"\n")
		})

		Convey("should add a timestamp if configured", func(c C) {
			at := "time"
			conf.AddTimestamp = &at
			m := mustNewMultiGrokRewriter(c, b, g, conf)
			runRewriter(c, m, []string{"hello"})
			So(b.String(), ShouldContainSubstring, `"any":"hello"`)
			So(b.String(), ShouldContainSubstring, `"time":"`)
		})

		Convey("should include default fields", func(c C) {
			conf.DefaultFields = map[string]interface{}{"foo": "bar", "baz": []int{1, 2, 3}}
			m := mustNewMultiGrokRewriter(c, b, g, conf)
			runRewriter(c, m, []string{"hello"})
			So(b.String(), ShouldEqual, `{"any":"hello","baz":[1,2,3],"foo":"bar"}`+"\n")
		})

		Convey("should capture input matching a complexish pattern", func(c C) {
			conf.MatchPatterns = []string{
				"%{NOTSPACE:perms} +%{INT:links:int} +%{NOTSPACE:user} +%{NOTSPACE:group} +%{INT:size:int} +%{LSTIMESTAMP:time} +%{GREEDYDATA:name}",
			}
			m := mustNewMultiGrokRewriter(c, b, g, conf)
			runRewriter(c, m, []string{"drwxr-xr-x   3 user  staff      102 Aug  2 13:31 vendor"})
			So(b.String(), ShouldEqual, `{"group":"staff","links":3,"name":"vendor","perms":"drwxr-xr-x","size":102,"time":"Aug  2 13:31","user":"user"}`+"\n")
		})

		Convey("should ignore input not matching a pattern", func(c C) {
			conf.MatchPatterns = []string{"%{INT:number}"}
			m := mustNewMultiGrokRewriter(c, b, g, conf)
			runRewriter(c, m, []string{"hello"})
			So(b.String(), ShouldEqual, "")
		})

		Convey("should handle typed patterns", func(c C) {
			conf.MatchPatterns = []string{"%{INT:number:int}"}
			m := mustNewMultiGrokRewriter(c, b, g, conf)
			runRewriter(c, m, []string{"123"})
			So(b.String(), ShouldEqual, `{"number":123}`+"\n")
		})

		Convey("with multilines", func() {
			Convey("should be able to match across lines with recursive fields", func(c C) {
				newLine := "\n"
				conf.MatchPatterns = []string{
					"%{WORD:hello} \n\t +%{WORD:world}\n\t %{GREEDYALL:numbers}",
				}
				conf.RecursiveFields = []config.GrokRewriterRecurseConfig{
					{
						Field:     "numbers",
						Delimiter: &newLine,
						Trim:      "\t ",
						Patterns:  []string{"%{INT:number:int}"},
					},
				}
				conf.BlacklistFields = []string{"numbers"}
				conf.Multiline.PrefixContinuations = []string{"\t"}
				m := mustNewMultiGrokRewriter(c, b, g, conf)

				runRewriter(c, m, []string{"hello ", "\t world", "\t 1", "\t 2", "\t 3"})

				So(b.String(), ShouldEqual, `{"hello":"hello","number":[1,2,3],"world":"world"}`+"\n")
			})
		})

		Convey("should support whitelisting", func(c C) {
			conf.MatchPatterns = []string{"%{WORD:one} %{WORD:two}"}
			conf.WhitelistFields = []string{"one"}
			m := mustNewMultiGrokRewriter(c, b, g, conf)
			runRewriter(c, m, []string{"hello there"})
			So(b.String(), ShouldEqual, `{"one":"hello"}`+"\n")
		})

		Convey("should support blacklisting", func(c C) {
			conf.MatchPatterns = []string{"%{WORD:one} %{WORD:two}"}
			conf.BlacklistFields = []string{"two"}
			m := mustNewMultiGrokRewriter(c, b, g, conf)
			runRewriter(c, m, []string{"hello there"})
			So(b.String(), ShouldEqual, `{"one":"hello"}`+"\n")
		})

		Convey("should support both blacklisting and whitelisting", func(c C) {
			conf.MatchPatterns = []string{"%{WORD:one} %{WORD:two} %{WORD:three}"}
			conf.WhitelistFields = []string{"one", "three"}
			conf.BlacklistFields = []string{"one", "two"}
			m := mustNewMultiGrokRewriter(c, b, g, conf)
			runRewriter(c, m, []string{"hello there world"})
			So(b.String(), ShouldEqual, `{"three":"world"}`+"\n")
		})
	})
}

func runRewriter(c C, r Rewriter, input []string) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		err := r.Run()
		wg.Done()
		c.So(err, ShouldBeNil)
	}()

	for _, line := range input {
		_, err := r.Write([]byte(line))
		c.So(err, ShouldBeNil)
		_, err = r.Write([]byte{'\n'})
		c.So(err, ShouldBeNil)
	}

	err := r.Close()
	So(err, ShouldBeNil)

	wg.Wait()
}

func mustNewMultiGrokRewriter(c C, out io.Writer, g *grok.Grok, conf *config.GrokRewriterConfig) *MultiGrokRewriter {
	m, err := NewMultiGrokRewriter(out, g, conf)
	c.So(err, ShouldBeNil)
	return m
}
