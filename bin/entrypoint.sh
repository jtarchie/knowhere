#!/bin/bash

set -eux

# Function to download database and restart server
download_and_restart() {
	if curl -s -q -o /var/osm/entries.db https://sqlite.knowhere.live/entries.db; then
		echo "Download complete. Restarting knowhere server..."
		# Find the process ID of knowhere server and send SIGKILL
		pkill -f "knowhere"
	else
		echo "Failed to download the database."
	fi
}

if [ ! -f /var/osm/entries.db ]; then
	download_and_restart &
fi

touch /var/osm/entries.db
/app/knowhere server --port 3000 --db /var/osm/entries.db --cors '*' --runtime-timeout 10s
