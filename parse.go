package main

import (
	"bufio"
	"os"
	"strings"
)

func readCmdLine() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	cmdLine, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	cmdLine = strings.TrimSuffix(cmdLine, "\n")
	return cmdLine, nil
}

func parseCmdLine(cmdLine string) ([]string, bool) {
	cmdLine, bg := removeBgSymb(cmdLine)
	cmd := strings.Split(cmdLine, " ")
	return cmd, bg
}

func removeBgSymb(cmdLine string) (string, bool) {
	if strings.HasSuffix(cmdLine, "&") {
		cmdLine = strings.TrimSuffix(cmdLine, "&")
		// Normalize "cmd &" and "cmd&".
		cmdLine = strings.TrimSuffix(cmdLine, " ")
		return cmdLine, true
	}
	return cmdLine, false
}
