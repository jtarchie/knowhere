package commands

import (
	"strings"
)

func connectionString(uri string) string {
	connectionString := uri + "?_query_only=true&immutable=true&mode=ro&_cache_size=5000&_busy_timeout=5000"

	if strings.Contains(uri, ".zst") {
		connectionString += "&vfs=zstd"
	}

	if !strings.HasPrefix(connectionString, "http") {
		connectionString = "file:" + connectionString
	}

	return connectionString
}
