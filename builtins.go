package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"
)

func isBuiltinCmd(cmd string) bool {
	switch cmd {
	case "kill", "jobs", "fg", "bg", "exit":
		return true
	}
	return false
}

func exit(cmd []string) error {
	os.Exit(1)
	return nil
}

func kill(cmd []string) error {
	sig := map[string]os.Signal{
		"-KILL": syscall.SIGKILL, "-9": syscall.SIGKILL,
		"-STOP": syscall.SIGTSTP, "-18": syscall.SIGTSTP,
		"-CONT": syscall.SIGCONT, "-19": syscall.SIGCONT,
		"-INT": syscall.SIGINT, "-2": syscall.SIGINT,
	}

	sigCmd := cmd[1]
	_, ok := sig[sigCmd]
	if !ok {
		return fmt.Errorf("kill: unknown argument")
	}

	pid, err := strconv.Atoi(cmd[2])
	if err != nil {
		return err
	}
	p, err := os.FindProcess(pid)
	if err != nil {
		return err
	}

	// Temp workaround SIGCONT issue.
	cmdStr := strings.Join(cmd, " ")
	if sigCmd == "-CONT" || sigCmd == "-19" {
		jobHandler(pid, contState, cmdStr)
	}

	p.Signal(sig[sigCmd])
	return nil
}

func lsJobs(cmd []string) error {
	for k, v := range jobsList {
		i := v[0]
		pid := k
		state := v[1]
		cmd := v[2]
		fmt.Printf("[%s] %d %s\t%s\n", i, pid, state, cmd)
	}
	return nil
}
