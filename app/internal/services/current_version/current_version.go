// Package current_version is for current version managing.

package current_version

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/ilya-mezentsev/versiond/app/internal/services/types"
	"github.com/ilya-mezentsev/versiond/app/internal/services/update"
)

const versionFilename = "versiond-current"

var (
	errCommandIsNotSpecified = errors.New("command-is-not-specified")
)

// Service manages demon current version value.
// It also implements update.VersionBoundCommand with the idea
// that Execute method will be called after version change.
type Service struct {
	cmd                    update.VersionBoundCommand
	currentVersionFilename string
	defaultVersion         string
}

func New(
	dir string,
	defaultVersion string,
) Service {

	return Service{
		currentVersionFilename: path.Join(dir, versionFilename),
		defaultVersion:         defaultVersion,
	}
}

func (s Service) MustGetCurrentVersion() string {
	dir := path.Dir(s.currentVersionFilename)
	if _, err := os.Stat(dir); errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir(dir, 0755)
		if err != nil {
			panic(fmt.Errorf("unable to create dir: %s, %v", dir, err))
		}
	}

	if _, err := os.Stat(s.currentVersionFilename); errors.Is(err, os.ErrNotExist) {
		err = s.save(s.defaultVersion)
		if err != nil {
			panic(fmt.Errorf("unable to save default version: %v", err))
		}

		return s.defaultVersion
	}

	bytes, err := os.ReadFile(s.currentVersionFilename)
	if err != nil {
		panic(fmt.Errorf("unable to read version from file: %s, %v", s.currentVersionFilename, err))
	}

	return strings.TrimSpace(string(bytes))
}

// Decorate wraps passed command so before its execution
// we can save new version as current.
func (s Service) Decorate(cmd update.VersionBoundCommand) Service {
	s.cmd = cmd
	return s
}

func (s Service) Execute(upd types.VersionUpd) error {
	err := s.save(upd.New)
	if err != nil {
		return err
	}

	if s.cmd == nil {
		return errCommandIsNotSpecified
	}

	return s.cmd.Execute(upd)
}

func (s Service) save(currentVersion string) error {
	return os.WriteFile(s.currentVersionFilename, []byte(currentVersion), 0755)
}
