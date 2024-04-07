// Package shell is for shell command execution

package shell

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/ilya-mezentsev/versiond/app/internal/services/types"
)

type Command struct {
	binaryPath string
}

func New(binaryPath string) Command {
	return Command{
		binaryPath: binaryPath,
	}
}

func (c Command) Execute(upd types.VersionUpd) error {
	if c.binaryPath == "" {
		return nil
	}

	cmd := exec.Command(c.binaryPath)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("CURRENT_VERSION=%s", upd.Current))
	cmd.Env = append(cmd.Env, fmt.Sprintf("NEW_VERSION=%s", upd.New))

	return cmd.Run()
}
