function cluster(query, radius) {
  const tree = rtree();

  return execute(`nwr[${query}](prefix=colorado)`).filter((entry) => {
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
  "#F0E442", // Yellow
  "#0072B2", // Blue
  "#D55E00", // Vermillion
  "#CC79A7", // Reddish Purple
  "#8DD3C7", // Light Blue-Green
  "#FDB462", // Soft Orange
  "#B3DE69", // Light Green
  "#FFED6F", // Light Yellow
  "#6A3D9A", // Deep Purple
  "#B15928", // Brownish-Orange
];

const keywords = [
  { query: "name=Costco", radius: 5 },
  { query: "amenity=cafe", radius: 1 },
  { query: "amenity=school", radius: 5 },
];

const clusters = keywords.map((keyword) => {
  return [keyword, cluster(keyword.query, keyword.radius)];
});

assert(clusters.length === keywords.length);

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

assertGeoJSON(payload);

return payload;
