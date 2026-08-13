//go:build linux

package integration

import (
	"os/exec"
	"syscall"
)

func syscallExec(command []string) error {
	path, err := exec.LookPath(command[0])
	if err != nil {
		return err
	}
	return syscall.Exec(path, command, syscall.Environ())
}
