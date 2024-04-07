package monitor

import (
	"errors"
	"testing"
	"time"

	monitorMocks "github.com/ilya-mezentsev/versiond/app/internal/services/monitor/mocks"
	"github.com/ilya-mezentsev/versiond/app/internal/services/types"
)

const (
	currentVersion = "v1"
	newestVersion  = "v2"
)

var (
	errNewestVersion = errors.New("error-newest-version")
	period           = time.Millisecond
)

func TestService_ListenOk(t *testing.T) {
	t.Parallel()

	expectedUpd := types.VersionUpd{
		Current: currentVersion,
		New:     newestVersion,
	}

	service := New(versionSourceMock(t, newestVersion, nil), currentVersion, period)

	updCh := make(chan types.VersionUpd)
	errCh := make(chan error)
	doneCh := make(chan struct{})

	go service.Listen(updCh, errCh, doneCh)

	upd := <-updCh
	if upd != expectedUpd {
		t.Errorf("upd diff: expected=%v, actual=%v", expectedUpd, upd)
	}

	doneCh <- struct{}{}

	close(updCh)
	close(errCh)
	close(doneCh)
}

func TestService_ListenOkSameVersions(t *testing.T) {
	t.Parallel()

	service := New(versionSourceMock(t, currentVersion, nil), currentVersion, period)

	updCh := make(chan types.VersionUpd)
	errCh := make(chan error)
	doneCh := make(chan struct{})

	go service.Listen(updCh, errCh, doneCh)

	// sleep here to make sure that update loop is executed at least once
	time.Sleep(5 * period)

	doneCh <- struct{}{}

	if service.currentVersion != currentVersion {
		t.Errorf("service.currentVersion != currentVersion, %s != %s", service.currentVersion, currentVersion)
	}

	close(updCh)
	close(errCh)
	close(doneCh)
}

func TestService_ListenErrGetNewestVersion(t *testing.T) {
	t.Parallel()

	service := New(versionSourceMock(t, "", errNewestVersion), currentVersion, period)

	updCh := make(chan types.VersionUpd)
	errCh := make(chan error)
	doneCh := make(chan struct{})

	go service.Listen(updCh, errCh, doneCh)

	err := <-errCh
	if err != errNewestVersion {
		t.Errorf("error diff: expected=%v, actual=%v", errNewestVersion, err)
	}

	doneCh <- struct{}{}

	close(updCh)
	close(errCh)
	close(doneCh)
}

func versionSourceMock(t *testing.T, version string, err error) VersionSource {
	mock := monitorMocks.NewMockVersionSource(t)
	mock.EXPECT().Fetch().Return(version, err)

	return mock
}
