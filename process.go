package stdjson

import (
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

type Process struct {
	Cmd    *exec.Cmd
	Stdout io.Writer
	Stderr io.Writer
}

func NewProcess(name string, arg ...string) *Process {
	p := &Process{
		Cmd: exec.Command(name, arg...),
	}

	// perform platform specific configuration for the process
	p.platformSpecificConfig()

	return p
}

func (p *Process) Run() error {
	p.Cmd.Stdin = os.Stdin
	p.Cmd.Stdout = p.Stdout
	p.Cmd.Stderr = p.Stderr

	if err := p.Cmd.Start(); err != nil {
		return err
	}

	shouldQuit := make(chan error)

	go func() {
		err := p.Cmd.Wait()
		shouldQuit <- err
	}()

	sig := make(chan os.Signal)

	signals := []os.Signal{syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGSTOP}
	signal.Notify(sig, signals...)
	defer signal.Reset(signals...)

	for {
		select {
		case err := <-shouldQuit:
			return err
		case s := <-sig:
			p.signal(s)
		}
	}
}

func (p *Process) WaitUntilRunning() {
	running := make(chan struct{}, 1)
	go func() {
		for {
			if p.Cmd != nil && p.Cmd.Process != nil && p.Cmd.Process.Pid != 0 {
				running <- struct{}{}
				break
			}
		}
	}()

	<-running
}

// Send signals to the process groups instead of a single process in an attempt to clean up
// after ourselves slightly more proper.
func (p *Process) signal(sig os.Signal) error {
	pid := p.Cmd.Process.Pid

	pgid, err := syscall.Getpgid(pid)
	if err != nil {
		return err
	}

	// use pgid, ref: http://unix.stackexchange.com/questions/14815/process-descendants
	if pgid == pid {
		pid = -1 * pid
	}

	target, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	return target.Signal(sig)
}
