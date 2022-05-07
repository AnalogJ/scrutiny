package shell

import (
	"bytes"
	"errors"
	"github.com/sirupsen/logrus"
	"io"
	"os/exec"
	"path"
	"strings"
)

type localShell struct{}

func (s *localShell) Command(logger *logrus.Entry, cmdName string, cmdArgs []string, workingDir string, environ []string) (string, error) {
	logger.Infof("Executing command: %s %s", cmdName, strings.Join(cmdArgs, " "))

	cmd := exec.Command(cmdName, cmdArgs...)
	var stdBuffer bytes.Buffer

	logWriters := []io.Writer{
		&stdBuffer,
	}
	if logger.Logger.Level == logrus.DebugLevel {
		logWriters = append(logWriters, logger.Logger.Out)
	}

	mw := io.MultiWriter(logWriters...)

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
