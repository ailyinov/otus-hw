package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
)

const defaultReturnCode = 1

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	c := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	for envVar, val := range env {
		if val.NeedRemove {
			err := os.Unsetenv(envVar)
			if err != nil {
				fmt.Println(err.Error())
			}
		} else {
			err := os.Setenv(envVar, val.Value)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}
	c.Env = os.Environ()

	err := c.Run()

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok { //nolint
			ws := exitError.Sys().(syscall.WaitStatus)
			returnCode = ws.ExitStatus()
		} else {
			log.Print(err.Error())
			returnCode = defaultReturnCode
		}
	} else {
		ws := c.ProcessState.Sys().(syscall.WaitStatus)
		returnCode = ws.ExitStatus()
	}

	return
}
