package testutil

import (
	"bytes"
	"github.com/Sirupsen/logrus"
	"io"
	"os"
	"sync"
	"testing"
)

type CapturingTest struct {
	*testing.T
	internalTest *testing.T
	Writer       *Writer
}

func NewCapturingTest(t *testing.T, setOutput bool, out io.Writer) *CapturingTest {
	if out == nil {
		out = os.Stdout
	}
	ct := &CapturingTest{
		t, t, NewWriter(out),
	}
	if setOutput {
		logrus.SetOutput(ct.Writer)
	}
	return ct
}

func (t *CapturingTest) Fail() {
	t.Writer.Dump()
	t.internalTest.Fail()
}

type Writer struct {
	buffer *bytes.Buffer
	out    io.Writer
	mu     sync.Mutex
}

func NewWriter(out io.Writer) *Writer {
	return &Writer{
		buffer: &bytes.Buffer{},
		out:    out,
	}
}

func (w *Writer) Write(b []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.buffer.Write(b)
}

func (w *Writer) Dump() (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.out.Write(w.buffer.Bytes())
}
