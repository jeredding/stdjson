package rewriter

import (
	"bytes"
	log "github.com/Sirupsen/logrus"
	"github.com/nkvoll/stdjson/config"
	"io"
	"strings"
	"time"
)

type MultilineBuffer struct {
	Lines               chan string
	delimiter           byte
	stripDelimiter      bool
	buffer              bytes.Buffer
	inputLines          chan *string
	written             chan struct{}
	prefixContinuations []string
	timeout             time.Duration
}

func NewMultilineBuffer(c *config.MultilineConfig) *MultilineBuffer {
	t := c.MustTimeoutDuration()
	if t == nil {
		timeout := 10 * time.Millisecond
		t = &timeout
	}

	var delimiter byte = '\n'
	if len(c.Delimiter) > 0 {
		delimiter = c.Delimiter[0]
	}

	return &MultilineBuffer{
		Lines:               make(chan string),
		delimiter:           delimiter,
		stripDelimiter:      c.StripLastDelimiter,
		buffer:              bytes.Buffer{},
		inputLines:          make(chan *string),
		written:             make(chan struct{}, 1024),
		prefixContinuations: c.PrefixContinuations,
		timeout:             *t,
	}
}

func (r *MultilineBuffer) Write(b []byte) (n int, err error) {
	n, err = r.buffer.Write(b)
	r.written <- struct{}{}
	return
}

func (r *MultilineBuffer) Close() error {
	close(r.written)
	log.Debugln("Closing multiline buffer")
	return nil
}

func (r *MultilineBuffer) Run() {
	defer func() {
		log.Debugln("Closing multiline output lines")
		close(r.Lines)
	}()
	currentLine := []string{}

	go func() {
		defer func() {
			log.Debugln("Stopped reading input from multiline buffer")
			close(r.inputLines)
		}()

		for range r.written {
			for {
				line, err := r.buffer.ReadString(r.delimiter)
				if err != nil {
					if err == io.EOF {
						break
					}
					return
				}

				if r.stripDelimiter {
					line = line[:len(line)-1]
				}

				r.inputLines <- &line
			}
		}
	}()

	emit := func() {
		if len(currentLine) > 0 {
			r.Lines <- strings.Join(currentLine, string(r.delimiter))
			currentLine = []string{}
		}
	}

readInput:
	for {
		select {
		case line := <-r.inputLines:
			if line == nil {
				emit()
				return
			}

			for _, prefix := range r.prefixContinuations {
				if strings.HasPrefix(*line, prefix) {
					currentLine = append(currentLine, *line)
					continue readInput
				}
			}

			emit()
			currentLine = append(currentLine, *line)
		case <-time.After(r.timeout):
			emit()
		}
	}
}
