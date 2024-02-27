#!/bin/bash

set -eux

mkdir -p .build/

states=("colorado" "massachusetts")

for state in "${states[@]}"; do
	filename="$state.osm.pbf"
	pushd .build/
	if [ ! -f "$filename" ]; then
		curl -o "$filename" "https://download.geofabrik.de/north-america/us/$state-latest.osm.pbf"
	fi
	popd

	go run -tags fts5 -race github.com/jtarchie/knowhere build \
		--osm ".build/$filename" \
		--db .build/entries.db \
		--prefix "$state"
done
