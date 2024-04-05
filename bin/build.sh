#!/bin/bash

set -eux

mkdir -p .build/

states=("california" "colorado" "florida" "massachusetts")

rm -Rf .build/entries.db*

for state in "${states[@]}"; do
	filename="$state.osm.pbf"
	pushd .build/
	if [ ! -f "$filename" ]; then
		curl -o "$filename" "https://download.geofabrik.de/north-america/us/$state-latest.osm.pbf"
	fi
	popd

	go run -tags fts5 github.com/jtarchie/knowhere build \
		--osm ".build/$filename" \
		--db .build/entries.db \
		--prefix "$state" --allowed-tags "name"
done

rclone copy .build/entries.db r2:knowhere-sqlite/ -P
