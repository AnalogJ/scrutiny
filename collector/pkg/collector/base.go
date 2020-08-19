package collector

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"
	"os/exec"
	"path"
)

type BaseCollector struct{}

func (c *BaseCollector) getJson(url string, target interface{}) error {

	r, err := httpClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func (c *BaseCollector) execCmd(cmdName string, cmdArgs []string, workingDir string, environ []string) (string, error) {

	cmd := exec.Command(cmdName, cmdArgs...)
	var stdBuffer bytes.Buffer
	mw := io.MultiWriter(os.Stdout, &stdBuffer)

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
