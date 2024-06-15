# Knowhere

This application allows you to search for locations of points of interests
(POIs) relative to each other.

Using the [`*.osm.pbf`](https://wiki.openstreetmap.org/wiki/PBF_Format) format,
rather than [XML](https://wiki.openstreetmap.org/wiki/OSM_XML), allows for the
transfer of smaller files for processing. The files (download from
[Geofabrik](https://download.geofabrik.de/)) are already compressed.

## Features

- Index OSM data into an SQLite data to help mass search of the Point of
  Interests.
- Javascript runtime for querying and manipulating the data.
- Ability to use a compressed database, which may require more CPU usage, but
  less resources overall.

## Usage

To run the tests:

```bash
brew bundle
task
```

### Build

Knowhere uses a read only database that it creates based on a list of
`*.osm.pbf` files. This conversion makes the data indexable for querying via
SQL.

When building the sqlite database, the OSM file are divided into separate
prefixed tables. This allows for isolation of regions, to make searching within
well defined boundaries easier. It may duplicate data across prefixed tables (as
OSM exports tend), but we have found the duplication to be minimal.

For example, to build each state:

```bash
./bin/build.sh ./bin/config/states.txt
```

This will create two artifacts `./build/entries.db` and `./build/entries.db.zst`
(zstd compressed version).

There will be tables with the prefix of each state:

- `alabama_entries`, `alabama_search`, and `alabama_rtree`
- ...
- `north_carolina_entries`, `north_carolina_search`, and `north_carolina_rtree`
- ...

The `entries` table contains (mostly) the original OSM data. The `search` table
contains FTS5 indexed data for full test search using `porter` tokenizer. The
`rtree` table optimized indexed table to entries within a bounding box.

### Deploy

This deployment is designed with [Fly](https://fly.io) in mind. When the
application starts for the first time, it will start the download of the sqlite
database, from a S3 file store. This requires a persistent file volume be added
to the container.

```bash
fly deploy
```

## TODO

- identify globbing over exact match
- add prefix list in runtime
