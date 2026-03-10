package auth

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// GetToken resolves a GitHub token using the following priority:
// 1. Explicit token (from --token flag)
// 2. gh CLI auth token
// 3. GITHUB_TOKEN environment variable
func GetToken(flagToken string) (string, error) {
	if flagToken != "" {
		return flagToken, nil
	}

	if token, err := getGHToken(); err == nil && token != "" {
		return token, nil
	}

	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		return token, nil
	}

	return "", fmt.Errorf("no GitHub token found\n\nTo authenticate, do one of:\n  1. Install and login with gh CLI: gh auth login\n  2. Set GITHUB_TOKEN environment variable\n  3. Pass --token flag")
}

func getGHToken() (string, error) {
	cmd := exec.Command("gh", "auth", "token")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
