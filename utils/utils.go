package utils

import (
	"fmt"
	"os"
)

// GetEnv parses an environment variable and returns it if available
func GetEnv(name string) (string, error) {
	v := os.Getenv(name)
	if v == "" {
		return "", fmt.Errorf("missing required environment variable %s", name)
	}

	return v, nil
}
