// Package command is for creating commands from config.
// At this moment only shell commands are supported.

package command

import "github.com/ilya-mezentsev/versiond/app/internal/services/command/shell"

func Shell(binaryPath string) shell.Command {
	return shell.New(binaryPath)
}
