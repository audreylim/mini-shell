package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
)

func execBuiltinCmd(cmd []string) {
	builtin := map[string]func([]string) error{
		"kill": kill,
		"jobs": lsJobs,
		"exit": exit,
	}
	err := builtin[cmd[0]](cmd)
	if err != nil {
		fmt.Println(err)
	}
}

func execFgCmd(cmd []string, sigStateChanged chan string) {
	cmdStr := strings.Join(cmd, " ")

	// TODO: Extract start process into common function.
	argv0, err := exec.LookPath(cmd[0])
	if err != nil {
		if cmd[0] != "" {
			fmt.Printf("Unknown command: %s\n", cmd[0])
		}
		// Don't execute new process with empty return. Will cause panic.
		sigPrompt <- struct{}{}
		return
	}
	var procAttr os.ProcAttr
	procAttr.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr}

	p, err := os.StartProcess(argv0, cmd, &procAttr)
	if err != nil {
		fmt.Printf("Start process %s, %s failed: %v", err, argv0, cmd)
	}

	for {
		sigChild := make(chan os.Signal)
		defer close(sigChild)
		// SIGCONT not receivable: https://github.com/golang/go/issues/8953
		// This causes some bugs. Eg. CONT signal not captured by handler means subsequent KILL or STOP signals will be ignored by this handler.
		signal.Notify(sigChild, syscall.SIGTSTP, syscall.SIGINT, syscall.SIGCONT, syscall.SIGKILL)
		defer signal.Stop(sigChild)

		var ws syscall.WaitStatus
		// Ignoring error. May return "no child processes" error. Eg. Sending Ctrl-c on `cat` command.
		wpid, _ := syscall.Wait4(p.Pid, &ws, syscall.WUNTRACED, nil)

		if ws.Exited() {
			break
		}
		if ws.Stopped() {
			jobHandler(wpid, runningState, cmdStr)
			jobHandler(wpid, suspendedState, cmdStr)
			// Return prompt when fg has become bg
			sigPrompt <- struct{}{}
		}
		//if ws.Continued() {
		//	state = contState
		//}
		if ws == 9 {
			jobHandler(wpid, killedState, cmdStr)
			break
		}
	}

	p.Wait()
	sigPrompt <- struct{}{}
}

func execBgCmd(cmd []string, sigStateChanged chan string) {
	cmdStr := strings.Join(cmd, " ")

	argv0, err := exec.LookPath(cmd[0])
	if err != nil {
		if cmd[0] != "" {
			fmt.Printf("Unknown command: %s\n", cmd[0])
		}
		sigPrompt <- struct{}{}
		return
	}
	var procAttr os.ProcAttr
	procAttr.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr}

	p, err := os.StartProcess(argv0, cmd, &procAttr)
	if err != nil {
		fmt.Printf("Start process %s, %s failed: %v", err, argv0, cmd)
	}
	jobHandler(p.Pid, runningState, cmdStr)
	sigPrompt <- struct{}{}

	//FIXME: Bg processes should not receive keyboard signals sent to fg process.

	for {
		sigChild := make(chan os.Signal)
		defer close(sigChild)
		signal.Notify(sigChild, syscall.SIGCHLD)
		defer signal.Stop(sigChild)

		var ws syscall.WaitStatus
		wpid, _ := syscall.Wait4(p.Pid, &ws, syscall.WUNTRACED, nil)

		if ws.Exited() {
			jobHandler(wpid, doneState, cmdStr)
			break
		}
		if ws.Stopped() {
			jobHandler(wpid, suspendedState, cmdStr)
			sigPrompt <- struct{}{}
		}
		//if ws.Continued() {
		//	state = contState
		//}
		if ws == 9 {
			jobHandler(wpid, killedState, cmdStr)
			break
		}
	}

	p.Wait()
	sigPrompt <- struct{}{}
}
