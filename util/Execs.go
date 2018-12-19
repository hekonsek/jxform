package util

import (
	"os"
	"os/exec"
	"strings"
)

type Execs struct {
}

func NewExecs() *Execs {
	return &Execs{}
}

func (execs *Execs) Run(command string, args ...string) ([]string, error) {
	out, err := exec.Command(command, args...).CombinedOutput()
	return strings.Split(string(out), "\n"), err
}

func (execs *Execs) Sout(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	err := cmd.Start()
	if err != nil {
		return err
	}

	err = cmd.Wait()
	return err
}
