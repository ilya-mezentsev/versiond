package update

import "github.com/ilya-mezentsev/versiond/app/internal/services/types"

type VersionBoundCommand interface {
	Execute(upd types.VersionUpd) error
}
