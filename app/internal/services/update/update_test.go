package update

import (
	"errors"
	"testing"

	"github.com/ilya-mezentsev/versiond/app/internal/services/types"
	updateMocks "github.com/ilya-mezentsev/versiond/app/internal/services/update/mocks"
)

var (
	upd = types.VersionUpd{
		Current: "v1",
		New:     "v2",
	}

	errBeforeCommand = errors.New("before-error")
	errMainCommand   = errors.New("main-error")
	errAfterCommand  = errors.New("after-error")
)

func TestService_Listen(t *testing.T) {
	tests := []struct {
		name          string
		beforeMock    VersionBoundCommand
		mainMock      VersionBoundCommand
		afterMock     VersionBoundCommand
		expectedError error
	}{
		{
			name:          "ok",
			beforeMock:    commandMock(t, nil),
			mainMock:      commandMock(t, nil),
			afterMock:     commandMock(t, nil),
			expectedError: nil,
		},

		{
			name:          "error from before command",
			beforeMock:    commandMock(t, errBeforeCommand),
			expectedError: errBeforeCommand,
		},

		{
			name:          "error from main command",
			beforeMock:    commandMock(t, nil),
			mainMock:      commandMock(t, errMainCommand),
			expectedError: errMainCommand,
		},

		{
			name:          "error from after command",
			beforeMock:    commandMock(t, nil),
			mainMock:      commandMock(t, nil),
			afterMock:     commandMock(t, errAfterCommand),
			expectedError: errAfterCommand,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := New(tt.beforeMock, tt.mainMock, tt.afterMock)

			updCh := make(chan types.VersionUpd)
			errsCh := make(chan error, 1)

			go service.Listen(updCh, errsCh)

			updCh <- upd
			if tt.expectedError != nil {
				if err := <-errsCh; err != tt.expectedError {
					t.Errorf("got an unexpected error: expected=%v, actual=%v", tt.expectedError, err)
				}
			}

			close(updCh)
			close(errsCh)
		})
	}
}

func commandMock(t *testing.T, err error) VersionBoundCommand {
	mock := updateMocks.NewMockVersionBoundCommand(t)
	mock.EXPECT().Execute(upd).Return(err)

	return mock
}
