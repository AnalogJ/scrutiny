package shell

import (
	"github.com/sirupsen/logrus"
)

// Create mock using:
// mockgen -source=collector/pkg/common/shell/interface.go -destination=collector/pkg/common/shell/mock/mock_shell.go
type Interface interface {
	Command(logger *logrus.Entry, cmdName string, cmdArgs []string, workingDir string, environ []string) (string, error)
}
