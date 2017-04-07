package util

import (
	"fmt"
	"os/exec"
	"os"
	"bytes"
)

func RunCommand(cwd string, env []string, cmd string, args ...string) (string, string, error) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	var c *exec.Cmd
	c = exec.Command(cmd, args...)
	c.Stdout = stdout
	c.Stderr = stderr
	c.Stdin = os.Stdin
	c.Dir = cwd
	c.Env = env

	err := c.Run()
	if err != nil {
		return stdout.String(), stderr.String(), fmt.Errorf("Error running command: %v\n%v", cmd, err)
	}

	if !c.ProcessState.Success() {
		return stdout.String(), stderr.String(), fmt.Errorf("Command exited in an error state: %v", cmd)
	}
	return stdout.String(), stderr.String(), err
}


