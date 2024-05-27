#!/bin/bash

set -eux

db_path=.build/entries.db

go run -tags fts5 github.com/jtarchie/knowhere build \
	--config "$1" \
	--db "$db_path" \
	--allowed-tags "name,amenity,shop,leisure,tourism"

go run github.com/SaveTheRbtz/zstd-seekable-format-go/cmd/zstdseek \
	-f "$db_path" \
	-o "$db_path".zst \
	-q 7 \
	-c 16:32:64

rclone copy "$db_path" r2:knowhere-sqlite/ -P
rclone copy "$db_path".zst r2:knowhere-sqlite/ -P

./bin/cleanup.sh
