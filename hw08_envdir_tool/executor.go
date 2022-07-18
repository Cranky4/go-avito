package main

import (
	"os"
	"os/exec"
	"strings"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	command := exec.Command(cmd[0], cmd[1:]...)

	for key, value := range env {
		if value.NeedRemove {
			os.Unsetenv(key)
		}
	}

	command.Env = append(command.Env, os.Environ()...)

	for key, value := range env {
		command.Env = append(command.Env, strings.Join([]string{key, value.Value}, "="))
	}

	command.Stdin = os.Stdin
	command.Stderr = os.Stderr
	command.Stdout = os.Stdout

	if err := command.Run(); err != nil {
		returnCode = 1
		return
	}

	returnCode = command.ProcessState.ExitCode()

	return
}
