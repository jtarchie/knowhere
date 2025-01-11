package commands

import (
	"fmt"
	"strings"
)

func connectionString(uri string) string {
	connectionString := fmt.Sprintf("file:%s?_query_only=true&immutable=true&mode=ro&_cache_size=5000&_busy_timeout=5000", uri)

	if strings.Contains(uri, ".zst") {
		connectionString += "&vfs=zstd"
	}

	return connectionString
}
