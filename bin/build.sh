#!/bin/bash

set -eux

mkdir -p .build/

states=("california" "colorado" "florida" "massachusetts")

db_path=.build/entries.db

rm -Rf "$db_path"

for state in "${states[@]}"; do
	filename="$state.osm.pbf"
	pushd .build/
	if [ ! -f "$filename" ]; then
		curl -o "$filename" "https://download.geofabrik.de/north-america/us/$state-latest.osm.pbf"
	fi
	popd

	go run -tags fts5 github.com/jtarchie/knowhere build \
		--osm ".build/$filename" \
		--db "$db_path" \
		--prefix "$state" --allowed-tags "name"
done

go run github.com/SaveTheRbtz/zstd-seekable-format-go/cmd/zstdseek \
	-f "$db_path" \
	-o "$db_path".zst \
	-q 7

rclone copy "$db_path" r2:knowhere-sqlite/ -P
rclone copy "$db_path".zst r2:knowhere-sqlite/ -P
