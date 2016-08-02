package stdjson

import "syscall"

func (p *Process) platformSpecificConfig() {
	if p.Cmd.SysProcAttr == nil {
		p.Cmd.SysProcAttr = &syscall.SysProcAttr{Pdeathsig: syscall.SIGKILL}
	} else {
		p.Cmd.SysProcAttr.Pdeathsig = syscall.SIGKILL
	}
}
