package monitor

type (
	VersionSource interface {
		Fetch() (string, error)
	}
)
