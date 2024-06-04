const colorPalette = [
  "#E69F00", // Orange
  "#56B4E9", // Sky Blue
  "#009E73", // Bluish Green
];

const keywords = [
  { query: "name=Costco", radius: 5 },
  { query: "amenity=cafe", radius: 1 },
  { query: "amenity=school", radius: 5 },
];

assert.stab("start");

keywords.forEach((keyword) => {
  keyword.results = geo.query(`nwr[${keyword.query}](prefix=colorado)`);
  assert.stab(`query ${keyword.query}`);
});

assert.stab("sort");
keywords.sort((a, b) => a.results.length - b.results.length);

assert.stab("cluster");
const neighbors = keywords[0].results.cluster(keywords[0].radius).map((
  entry,
) => [entry]);

assert.stab("closeby");
keywords.slice(1).forEach((keyword) => {
  assert.stab(`tree ${keyword.query}`);
  const tree = keyword.results.asTree(keyword.radius)

  assert.stab(`neighbor ${keyword.query}`);
  neighbors.forEach((entries) => {
    const extended = entries[0].bbox().extend(keywords[0].radius);

    const nearby = tree.nearby(extended, 1);
    if (nearby.length === 1) {
      entries.push(nearby[0]);
    }
  });
});

assert.stab("payload");
const payload = {
  type: "FeatureCollection",
  features: neighbors.flatMap((entries) => {
    if (entries.length !== keywords.length) {
      return;
    }

    return entries.flatMap((entry, index) => {
      const color = colorPalette[index % entries.length];

      const feature = entry.asFeature({
        "marker-color": color,
        index: index,
      });

      return feature;
    });
  }),
};

assert.stab("assert");
assert.geoJSON(payload);

assert.stab("return");
return payload;
