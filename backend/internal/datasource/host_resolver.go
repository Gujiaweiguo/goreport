package datasource

import (
	"os"
	"strings"
)

const datasourceLocalhostAliasEnv = "DATASOURCE_LOCALHOST_ALIAS"

func ResolveHost(host string) string {
	trimmedHost := strings.TrimSpace(host)
	if trimmedHost == "" {
		return trimmedHost
	}

	alias := strings.TrimSpace(os.Getenv(datasourceLocalhostAliasEnv))
	if alias == "" {
		return trimmedHost
	}

	switch strings.ToLower(trimmedHost) {
	case "localhost", "127.0.0.1", "::1":
		return alias
	default:
		return trimmedHost
	}
}
