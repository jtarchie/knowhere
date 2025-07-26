#!/bin/bash

set -eux

filename="entries.db.zst"

if [ ! -f /var/osm/$filename ]; then
  # Start the server in the background
  /app/knowhere server --port 3000 --log-level debug --db /var/osm/$filename --cors "${CORS_DOMAIN:-*}" --allow-cidr "${ALLOW_CIDR:-0.0.0.0/0}" --runtime-timeout 30s &
  SERVER_PID=$!
  
  # Download the database
  curl -q --progress-bar -o /var/osm/entries.db.zst https://sqlite.knowhere.live/entries.db.zst
  
  # Kill the old server
  kill $SERVER_PID
  wait $SERVER_PID 2>/dev/null || true
  
  echo "Download complete, restarting server with updated database..."
fi

# Start the server (either first time with existing DB, or restart after download)
exec /app/knowhere server --port 3000 --log-level debug --db /var/osm/$filename --cors "${CORS_DOMAIN:-*}" --allow-cidr "${ALLOW_CIDR:-0.0.0.0/0}" --runtime-timeout 30s