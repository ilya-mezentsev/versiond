// Package services is for demon services initialization

package services

import (
	"github.com/ilya-mezentsev/versiond/app/internal/services/command"
	"github.com/ilya-mezentsev/versiond/app/internal/services/current_version"
	"github.com/ilya-mezentsev/versiond/app/internal/services/monitor"
	"github.com/ilya-mezentsev/versiond/app/internal/services/request"
	"github.com/ilya-mezentsev/versiond/app/internal/services/types"
	"github.com/ilya-mezentsev/versiond/app/internal/services/update"
)

type Services struct {
	update  update.Service
	monitor monitor.Service
}

func New(config Config) Services {
	currentVersion := current_version.New(
		config.Cache.Dir,
		config.DefaultVersion,
	)

	return Services{
		monitor: monitor.New(
			request.New(
				config.Version.Source.Address,
				config.Version.Source.Timout,
			),
			currentVersion.MustGetCurrentVersion(),
			config.Version.Monitor.Period,
		),

		update: update.New(
			command.Shell(config.OnChangeCmds.Before),
			command.Shell(config.OnChangeCmds.Main),
			currentVersion.Decorate(command.Shell(config.OnChangeCmds.After)),
		),
	}
}

func (ss Services) Run(
	done <-chan struct{},
	errs chan<- error,
) {

	updatesChan := make(chan types.VersionUpd)
	monitorErrorsChan := make(chan error)
	updateErrorsChan := make(chan error)
	monitorDoneChan := make(chan struct{})

	go ss.monitor.Listen(
		updatesChan,
		monitorErrorsChan,
		monitorDoneChan,
	)
	go func() {
		for monitorErr := range monitorErrorsChan {
			errs <- monitorErr
		}
	}()

	go ss.update.Listen(
		updatesChan,
		updateErrorsChan,
	)
	go func() {
		for updateErr := range updateErrorsChan {
			errs <- updateErr
		}
	}()

	<-done

	monitorDoneChan <- struct{}{}
	close(updatesChan)
	close(monitorErrorsChan)
	close(updateErrorsChan)
	close(monitorDoneChan)
	close(errs)
}
