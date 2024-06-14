#!/bin/bash

set -euo pipefail # Improved error handling

# Default values
db_path=".build/entries.db"
enable_rclone=true
enable_cleanup=true

config_path=$1
shift

# Parse command-line arguments
while [[ $# -gt 0 ]]; do
	key="$1"
	case $key in
	--db-path)
		db_path="$2"
		shift 2
		;;
	--no-rclone)
		enable_rclone=false
		shift
		;;
	--no-cleanup)
		enable_cleanup=false
		shift
		;;
	*)
		echo "Unknown option: $key" >&2
		exit 1
		;;
	esac
done

# Building the SQLite database
go run -tags fts5 github.com/jtarchie/knowhere build \
	--config "$config_path" \
	--db "$db_path" \
	--allowed-tags "name,amenity,shop,leisure,tourism,boundary,admin_level,waterway,border_type"

# Compressing the database
go run github.com/SaveTheRbtz/zstd-seekable-format-go/cmd/zstdseek \
	-f "$db_path" \
	-o "$db_path".zst \
	-q 7 \
	-c 16:32:64

# Rclone copy (if enabled)
if $enable_rclone; then
	rclone copy "$db_path" r2:knowhere-sqlite/ -P
	rclone copy "$db_path".zst r2:knowhere-sqlite/ -P
fi

# Cleanup (if enabled)
if $enable_cleanup; then
	./bin/cleanup.sh
fi
