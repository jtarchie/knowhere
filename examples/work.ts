/// <reference path="../docs/global.d.ts" />

const homes = [
  // lat, lng, address
  "331 Heather Hill Dr, Gibsonia, PA 15044",
  "200 Pine Mountain Ln, McCandless, PA 15090",
].flatMap((address) => {
  const entries = query.fromAddress(address);
  assert.eq(
    entries.length >= 1,
    `expected one address match ${entries.length}`,
  );
  return entries.map((
    entry,
  ): [number, number, string] => [entry.minLat, entry.minLon, entry.tags.name]);
});

const impacts = query.execute(
  `nwr[name=~"Western Psychiatric"](prefix=pennsylvania)`,
).map((entry): [number, number, string] => {
  return [entry.minLat, entry.minLon, entry.tags.name];
});

let features = homes.map((coords, index) => {
  const point = geo.asPoint(coords[0], coords[1]);

  return point.asFeature({
    "marker-color": colors.pick(index),
    "title": coords[2],
  });
});

features = features.concat(impacts.map((coords, index) => {
  const point = geo.asPoint(coords[0], coords[1]);

  return point.asFeature({
    "marker-color": colors.pick(index + homes.length + 1),
    "title": coords[2],
    "isochrone": true,
  });
}));

const payload = {
  type: "FeatureCollection",
  features: features,
};

assert.geoJSON(payload);

export { payload };
