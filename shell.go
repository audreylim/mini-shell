package main

import (
	"fmt"
)

var (
	sigStateChanged = make(chan string)
	sigPrompt       = make(chan struct{})

	// Not thread safe.
	jobsList = make(map[int][]string)
)

const (
	runningState    = "running"
	suspendedState  = "suspended"
	contState       = "continue"
	killedState     = "killed"
	terminatedState = "terminated"
	doneState       = "done"
)

const prettyPrompt = "\033[35m\u2764\033[m \033[36m\u2764\033[m \033[37m\u2764\033  "

func main() {
prompt:
	fmt.Printf(prettyPrompt)

	cmdLine, err := readCmdLine()
	if err != nil {
		fmt.Println(err)
	}
	cmd, bg := parseCmdLine(cmdLine)

	if isBuiltinCmd(cmd[0]) {
		execBuiltinCmd(cmd)
	} else {
		if bg {
			go execBgCmd(cmd, sigStateChanged)
			<-sigPrompt
		} else {
			go execFgCmd(cmd, sigStateChanged)
			<-sigPrompt
		}
	}

	goto prompt
}
