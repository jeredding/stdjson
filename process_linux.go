package stdjson

import "syscall"

func (p *Process) platformSpecificConfig() {
	p.Cmd.SysProcAttr = syscall.SysProcAttr{Pdeathsig: syscall.SIGKILL}
}
