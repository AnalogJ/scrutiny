package errors

import (
	"fmt"
)

// Raised when config file is missing
type ConfigFileMissingError string

func (str ConfigFileMissingError) Error() string {
	return fmt.Sprintf("ConfigFileMissingError: %q", string(str))
}

// Raised when the config file doesnt match schema
type ConfigValidationError string

func (str ConfigValidationError) Error() string {
	return fmt.Sprintf("ConfigValidationError: %q", string(str))
}

// Raised when a dependency (like smartd or ssh-agent) is missing
type DependencyMissingError string

func (str DependencyMissingError) Error() string {
	return fmt.Sprintf("DependencyMissingError: %q", string(str))
}

// Raised when the notification system is incorrectly configured
type NotificationValidationError string

func (str NotificationValidationError) Error() string {
	return fmt.Sprintf("NotificationValidationError: %q", string(str))
}
