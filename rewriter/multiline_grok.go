package rewriter

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/nkvoll/stdjson/config"
	"github.com/vjeantet/grok"
	"io"
	"reflect"
	"strings"
	"time"
)

type MultiGrokRewriter struct {
	out             io.Writer
	multiline       *MultilineBuffer
	g               *grok.Grok
	addTimestamp    *string
	matchPatterns   []string
	recursiveFields []config.GrokRewriterRecurseConfig
	whitelistFields []string
	blacklistFields []string
}

func NewMultiGrokRewriter(
	out io.Writer,
	g *grok.Grok,
	c *config.GrokRewriterConfig) *MultiGrokRewriter {
	r := &MultiGrokRewriter{
		out:             out,
		multiline:       NewMultilineBuffer(&c.Multiline),
		g:               g,
		addTimestamp:    c.AddTimestamp,
		matchPatterns:   c.MatchPatterns,
		recursiveFields: c.RecursiveFields,
		whitelistFields: c.WhitelistFields,
		blacklistFields: c.BlacklistFields,
	}
	return r
}

func (r *MultiGrokRewriter) Write(b []byte) (n int, err error) {
	n, err = r.multiline.Write(b)
	return
}

func (r *MultiGrokRewriter) Close() error {
	return r.multiline.Close()
}

func (r *MultiGrokRewriter) Run() error {
	go r.multiline.Run()

	for {
		select {
		case line, ok := <-r.multiline.Lines:
			if !ok {
				return nil
			}

			res := make(map[string]interface{})

			r.doGrok(res, line, r.matchPatterns)

			if len(r.whitelistFields) != 0 {
				whitelistedResult := make(map[string]interface{})
				for _, fieldName := range r.whitelistFields {
					if v, ok := res[fieldName]; ok {
						whitelistedResult[fieldName] = v
					}
				}
				res = whitelistedResult
			}

			for _, fieldName := range r.blacklistFields {
				delete(res, fieldName)
			}

			if len(res) == 0 {
				// found no parsable result, makes no sense to use this
				continue
			}

			if r.addTimestamp != nil {
				res[*r.addTimestamp] = time.Now().UTC().Format(time.RFC3339)
			}

			d, err := json.Marshal(&res)
			if err != nil {
				log.Errorln("unable to marshal to json", err)
			} else {
				d = append(d, '\n')
				_, err := r.out.Write(d)
				if err != nil {
					log.Errorln("unable to write to out", err)
				}
			}
		}
	}
}

func (r *MultiGrokRewriter) doGrok(target map[string]interface{}, line string, patterns []string) {
	for _, pattern := range patterns {
		log.Debugln("Trying to match pattern", pattern, "against line", line)

		v, err := r.g.ParseTyped(pattern, line)
		if err != nil {
			log.Errorln("error while attempting to parse", err)
			continue
		}

		if len(v) == 0 {
			continue
		}

		// TODO: proper map merging, nested keys, nested arrays w/ subkeys etc
		for key, value := range v {
			if current, ok := target[key]; !ok {
				target[key] = value
			} else {
				cv := reflect.ValueOf(current)

				switch cv.Kind() {
				case reflect.Slice | reflect.Array:
					// merge with existing slice / array
					ns := []interface{}{}
					for i := 0; i < cv.Len(); i++ {
						ns = append(ns, cv.Index(i).Interface())
					}
					target[key] = append(ns, value)
				default:
					// could be map, just create a new list and start appending
					ns := []interface{}{}
					ns = append(ns, cv.Interface())
					target[key] = append(ns, value)
				}
			}

			if vstr, ok := value.(string); ok {
				for _, field := range r.recursiveFields {
					if field.Field == key {
						vstr = strings.Trim(vstr, field.Trim)

						vlines := []string{vstr}

						if field.Delimiter != nil {
							vlines = strings.Split(vstr, *field.Delimiter)
						}

						for _, vline := range vlines {
							r.doGrok(target, vline, field.Patterns)
						}
					}
				}
			}
		}

		break
	}
}
