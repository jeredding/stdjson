package rewriter

import (
	"io"
	"sync"
)

type Noop struct {
	out     io.Writer
	running sync.WaitGroup
}

func NewNoop(out io.Writer) *Noop {
	r := &Noop{
		out: out,
	}

	r.running.Add(1)

	return r
}

func (r *Noop) Write(b []byte) (int, error) {
	return r.out.Write(b)
}

func (r *Noop) Run() error {
	r.running.Wait()
	return nil
}

func (r *Noop) Close() error {
	r.running.Done()
	return nil
}
