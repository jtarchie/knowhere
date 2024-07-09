#!/bin/bash

set -eux

filename="entries.db.zst"

# Function to download database and restart server
download_and_restart() {
	if curl -s -q -o /var/osm/$filename https://sqlite.knowhere.live/$filename; then
		echo "Download complete. Restarting knowhere server..."
		# Find the process ID of knowhere server and send SIGKILL
		pkill -f "knowhere"
	else
		echo "Failed to download the database."
	fi
}

if [ ! -f /var/osm/$filename ]; then
	echo "Downloading the database"
	download_and_restart &
fi

touch /var/osm/$filename
/app/knowhere server --port 3000 --db /var/osm/$filename --cors "${CORS_DOMAIN:-*}" --allow-cidr "${ALLOW_CIDR:-0.0.0.0/0}" --runtime-timeout 10s
