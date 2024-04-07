#!/bin/bash

set -eux

mkdir -p .build/

states=("alabama" "alaska" "arizona" "arkansas" "california" "colorado" "connecticut" "delaware" "florida" "georgia" "hawaii" "idaho" "illinois" "indiana" "iowa" "kansas" "kentucky" "louisiana" "maine" "maryland" "massachusetts" "michigan" "minnesota" "mississippi" "missouri" "montana" "nebraska" "nevada" "new-hampshire" "new-jersey" "new-mexico" "new-york" "north-carolina" "north-dakota" "ohio" "oklahoma" "oregon" "pennsylvania" "rhode-island" "south-carolina" "south-dakota" "tennessee" "texas" "utah" "vermont" "virginia" "washington" "west-virginia" "wisconsin" "wyoming")

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
