function cluster(query, radius) {
  const tree = geo.rtree();

  return geo.query(`nwr[${query}](prefix=colorado)`).filter((entry) => {
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
  { query: "name=Costco", radius: 5000 },
  { query: "amenity=cafe", radius: 1000 },
  { query: "amenity=school", radius: 5000 },
];

const clusters = keywords.map((keyword) => {
  return [keyword, cluster(keyword.query, keyword.radius)];
});

assert.eq(
  clusters.length === keywords.length,
  "expected same number of clusters",
);

const payload = {
  type: "FeatureCollection",
  features: clusters.flatMap((cluster, index) => {
    const color = colorPalette[index % clusters.length];

    return cluster[1].map((entry) => {
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
