package daemon

import (
	"errors"
	"net/http"
)

func Ping() error {
	resp, err := http.Get("http://localhost:7777/health")
	if err != nil {
		return errors.New("daemon is down")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("daemon unhealthy")
	}

	return nil
}
