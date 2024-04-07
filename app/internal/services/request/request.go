// Package request is for requesting the newest version from some backend.

package request

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type (
	Service struct {
		addr   string
		client http.Client
	}

	response struct {
		NewestVersion string `json:"newest_version"`
	}
)

func New(
	addr string,
	timeout time.Duration,
) Service {

	return Service{
		addr: addr,
		client: http.Client{
			Timeout: timeout,
		},
	}
}

func (s Service) Fetch() (string, error) {
	r, err := s.client.Get(s.addr)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = r.Body.Close()
	}()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return "", err
	}

	var versionResponse response
	err = json.Unmarshal(body, &versionResponse)
	if err != nil {
		return "", err
	}

	return versionResponse.NewestVersion, nil
}
