#!/bin/bash

set -eux

curl -vvv https://knowhere.fly.dev
fly ssh console --command "rm -Rf /var/osm/entries.db"
fly app restart