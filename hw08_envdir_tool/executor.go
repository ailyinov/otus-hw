package main

import (
	"log"
	"os"
	"os/exec"
	"syscall"
)

const defaultReturnCode = 1

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	c := exec.Command(cmd[0], cmd[1:]...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Env = os.Environ()
	for envVar, val := range env {
		if val.NeedRemove {
			c.Env = append(c.Env, envVar+"=")
		} else {
			c.Env = append(c.Env, envVar+"="+val.Value)
		}
	}

	err := c.Run()

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			ws := exitError.Sys().(syscall.WaitStatus)
			returnCode = ws.ExitStatus()
		} else {
			log.Printf(err.Error())
			returnCode = defaultReturnCode
		}
	} else {
		ws := c.ProcessState.Sys().(syscall.WaitStatus)
		returnCode = ws.ExitStatus()
	}

	return
}
