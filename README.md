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

## Usage

To run the tests:

```bash
brew bundle
task
```

variables: RAILS_ENV: production SECRET_KEY_BASE:
4b737bf9efeb7bd42ae7650ca823ecd5bb271588cc322eb8d670293a23d504939bacfa999598f44fd5d216a15aee694d1ca871dbdfc9c18a787b3f5805f141d6

secrets: DATABASE_URL: from_cfn:
${COPILOT_APPLICATION_NAME}-${COPILOT_ENVIRONMENT_NAME}-dbAuroraSecret
