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
  "#44AA99", // Teal  
  "#117733", // Dark Green
  "#999933", // Olive Green
  "#AA4499", // Purple
  "#DDCC77", // Light Tan 
  "#882255", // Dark Red
  "#332288", // Dark Blue
];

const keywords = [
  { query: "name=Costco", radius: 5000 },
  { query: "amenity=cafe", radius: 1000 },
  { query: "amenity=school", radius: 5000 },
];

assert.stab("start");

keywords.forEach((keyword) => {
  keyword.results = geo.query(`nwr[${keyword.query}](prefix=colorado)`);
  assert.stab(`query ${keyword.query}`);
});

assert.stab("sort");
keywords.sort((a, b) => a.results.length - b.results.length);

assert.stab("cluster");
const neighbors = keywords[0].results.cluster(500).map((
  entry,
) => geo.asResults(entry));

assert.stab("closeby");
keywords.slice(1).forEach((keyword) => {
  assert.stab(`tree ${keyword.query}`);
  const tree = keyword.results.asTree(keyword.radius);

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
  features: neighbors.flatMap((entries, index) => {
    if (entries.length !== keywords.length) {
      return;
    }

    const features = entries.flatMap((entry, index) => {
      const color = colorPalette[index % entries.length];

      const feature = entry.asFeature({
        "marker-color": color,
        index: index,
      });

      return feature;
    });

    const bounds = geo.asBounds(
      ...entries.map((entry, index) =>
        entry.bbox().extend(keywords[index].radius)
      ),
    );

    return features.concat(
      [
        bounds.asFeature({
          "fill": colorPalette[index % neighbors.length],
          "fill-opacity": 0.5,
        }),
      ],
    );
  }),
};

assert.stab("assert");
assert.geoJSON(payload);

assert.stab("return");
return payload;
