#!/bin/bash

set -eux

filename="entries.db.zst"

/app/knowhere server --port 3000 --log-level debug --db /var/osm/$filename --cors "${CORS_DOMAIN:-*}" --allow-cidr "${ALLOW_CIDR:-0.0.0.0/0}" --runtime-timeout 30s
