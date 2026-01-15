package remote

import (
	"fmt"
	"strconv"
	"strings"
)

func parsePid(out string) (int, error) {
	t := strings.TrimSpace(out)
	pid, err := strconv.Atoi(t)
	if err != nil {
		return 0, fmt.Errorf("invalid pid output: %q", t)
	}
	return pid, nil
}
