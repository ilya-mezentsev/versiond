// Package monitor is for version updates monitoring.

package monitor

import (
	"time"

	"github.com/ilya-mezentsev/versiond/app/internal/services/types"
)

type Service struct {
	newestVersionSource VersionSource
	currentVersion      string
	period              time.Duration
}

func New(
	newestVersionSource VersionSource,
	currentVersion string,
	period time.Duration,
) Service {

	return Service{
		newestVersionSource: newestVersionSource,
		currentVersion:      currentVersion,
		period:              period,
	}
}

func (s Service) Listen(
	updates chan<- types.VersionUpd,
	errs chan<- error,
	done <-chan struct{},
) {

	ticker := time.NewTicker(s.period)
	for {
		select {
		case <-ticker.C:
			s.checkUpdates(updates, errs)
		case <-done:
			ticker.Stop()
			return
		}
	}
}

func (s *Service) checkUpdates(
	update chan<- types.VersionUpd,
	errs chan<- error,
) {

	newestVersion, err := s.newestVersionSource.Fetch()
	if err != nil {
		errs <- err
		return
	}

	if newestVersion != s.currentVersion {
		update <- types.VersionUpd{
			Current: s.currentVersion,
			New:     newestVersion,
		}

		s.currentVersion = newestVersion
	}
}
