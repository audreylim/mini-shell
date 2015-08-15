package main

import (
	"fmt"
	"strconv"
)

// FIXME: There must be a better way to organise this.
func jobHandler(pid int, state string, cmd string) {
	var job string
	switch state {
	case runningState:
		_, ok := jobsList[pid]
		if !ok {
			i := strconv.Itoa(len(jobsList) + 1)

			job = fmt.Sprintf("[%s] %d %s\n", i, pid, state)
			jobInfo := []string{i, state, cmd}
			jobsList[pid] = jobInfo
		}
		sigPrompt <- struct{}{}
	case suspendedState:
		jobInfo := jobsList[pid]

		i := jobInfo[0]
		cmdLine := jobInfo[2]
		job = fmt.Sprintf("[%s] %d %s\t%s\n", i, pid, state, cmdLine)
		jobInfo = []string{i, state, cmd}
		jobsList[pid] = jobInfo

	case contState:
		jobInfo := jobsList[pid]

		i := jobInfo[0]
		cmdLine := jobInfo[2]
		job = fmt.Sprintf("[%s] %d %s\t%s\n", i, pid, state, cmdLine)
		jobInfo = []string{i, runningState, cmdLine}
		jobsList[pid] = jobInfo

	case killedState:
		jobInfo := jobsList[pid]
		i := jobInfo[0]
		cmdLine := jobInfo[2]
		job = fmt.Sprintf("[%s] %d %s\t%s\n", i, pid, state, cmdLine)
		delete(jobsList, pid)
	case terminatedState:
		jobInfo := jobsList[pid]
		i := jobInfo[0]
		cmdLine := jobInfo[2]
		job = fmt.Sprintf("[%s] %d %s\t%s\n", i, pid, state, cmdLine)
		delete(jobsList, pid)
	case doneState:
		jobInfo := jobsList[pid]
		i := jobInfo[0]
		cmdLine := jobInfo[2]
		job = fmt.Sprintf("[%s] %d %s\t%s\n", i, pid, state, cmdLine)
		delete(jobsList, pid)
	default:
	}
	fmt.Printf(job)
	return
}
