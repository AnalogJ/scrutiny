package common

import (
	"bytes"
	"errors"
	"github.com/sirupsen/logrus"
	"io"
	"os/exec"
	"path"
	"strings"
)

func ExecCmd(logger *logrus.Entry, cmdName string, cmdArgs []string, workingDir string, environ []string) (string, error) {
	logger.Infof("Executing command: %s %s", cmdName, strings.Join(cmdArgs, " "))

	cmd := exec.Command(cmdName, cmdArgs...)
	var stdBuffer bytes.Buffer
	mw := io.MultiWriter(logger.Logger.Out, &stdBuffer)

	cmd.Stdout = mw
	cmd.Stderr = mw

	if environ != nil {
		cmd.Env = environ
	}
	if workingDir != "" && path.IsAbs(workingDir) {
		cmd.Dir = workingDir
	} else if workingDir != "" {
		return "", errors.New("Working Directory must be an absolute path")
	}

	err := cmd.Run()
	return stdBuffer.String(), err

}
