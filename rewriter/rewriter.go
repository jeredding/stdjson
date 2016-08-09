package rewriter

import (
	"errors"
	"github.com/nkvoll/stdjson/config"
	"github.com/vjeantet/grok"
	"io"
)

type Rewriter interface {
	io.WriteCloser
	Run() error
}

func RewriterForStreamConfig(sc *config.StreamConfig, out io.Writer) (Rewriter, error) {
	if sc.Rewriter.GrokRewriter != nil {
		namedOnly := true
		if sc.Rewriter.GrokRewriter.NamedOnly != nil {
			namedOnly = *sc.Rewriter.GrokRewriter.NamedOnly
		}

		g, err := grok.NewWithConfig(&grok.Config{
			NamedCapturesOnly: namedOnly,
		})

		if err != nil {
			return nil, err
		}

		c := sc.Rewriter.GrokRewriter
		for _, ep := range c.ExtraPatterns {
			if ep.Inline != nil {
				if err := g.AddPattern(ep.Inline.Name, ep.Inline.Pattern); err != nil {
					return nil, err
				}
			}
		}

		return NewMultiGrokRewriter(
			out, g, c,
		)
	}

	return nil, errors.New("no supported rewriters found")
}
