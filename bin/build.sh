#!/bin/bash

set -eux

mkdir -p .build/

db_path=.build/entries.db

rm -Rf "$db_path"*

while read -r url; do
	filename="$(basename "$url")"
	country_code="$(basename "$url" | awk -F'[/-]' '{print $(NF-1)}')"

	pushd .build/
	if [ ! -f "$filename" ]; then
		curl -o "$filename" "$url"
	fi
	popd

	go run -tags fts5 github.com/jtarchie/knowhere build \
		--osm ".build/$filename" \
		--db "$db_path" \
		--prefix "$country_code" --allowed-tags "name"
done <"$1"

go run github.com/SaveTheRbtz/zstd-seekable-format-go/cmd/zstdseek \
	-f "$db_path" \
	-o "$db_path".zst \
	-q 7

# rclone copy "$db_path" r2:knowhere-sqlite/ -P
# rclone copy "$db_path".zst r2:knowhere-sqlite/ -P

# ./bin/cleanup.sh
