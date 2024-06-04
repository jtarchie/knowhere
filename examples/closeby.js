function cluster(results, radius) {
  const tree = geo.rtree();

  return results.filter((entry) => {
    const extended = entry.bbox().extend(radius);

    if (tree.within(extended)) {
      return false;
    }

    tree.insert(extended);
    return true;
  });
}

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

keywords.forEach((keyword) => {
  keyword.results = geo.query(`nwr[${keyword.query}](prefix=colorado)`);
});

keywords.sort((a, b) => a.results.length - b.results.length);

const neighbors = cluster(keywords[0].results, keywords[0].radius).map((
  entry,
) => [entry]);

keywords.slice(1).forEach((keyword) => {
  const tree = geo.rtree();

  keyword.results.forEach((entry) => {
    const extended = entry.bbox().extend(keyword.radius);

    tree.insert(extended, entry);
  });

  neighbors.forEach((entries) => {
    const extended = entries[0].bbox().extend(keywords[0].radius);

    const nearby = tree.nearby(extended, 1);
    if (nearby.length === 1) {
      entries.push(nearby[0]);
    }
  });
});

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

assert.geoJSON(payload);

return payload;
