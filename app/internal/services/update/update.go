// Package update is for performing action after version update.

package update

import "github.com/ilya-mezentsev/versiond/app/internal/services/types"

type Service struct {
	// beforeCommand - command executed before main
	beforeCommand VersionBoundCommand

	// mainCommand - command that supposed to be "main" for reaction on version update
	mainCommand VersionBoundCommand

	// afterCommand - command executed after main
	afterCommand VersionBoundCommand
}

func New(
	beforeCommand VersionBoundCommand,
	mainCommand VersionBoundCommand,
	afterCommand VersionBoundCommand,
) Service {

	return Service{
		beforeCommand: beforeCommand,
		mainCommand:   mainCommand,
		afterCommand:  afterCommand,
	}
}

func (s Service) Listen(
	updates <-chan types.VersionUpd,
	errs chan<- error,
) {

	for upd := range updates {
		err := s.processUpdate(upd)
		if err != nil {
			errs <- err
		}
	}
}

func (s Service) processUpdate(upd types.VersionUpd) error {
	if err := s.beforeCommand.Execute(upd); err != nil {
		return err
	}

	if err := s.mainCommand.Execute(upd); err != nil {
		return err
	}

	return s.afterCommand.Execute(upd)
}
