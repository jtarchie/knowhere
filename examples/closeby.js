function zillowURL(bounds) {
  // https://www.zillow.com/homes/for_sale/?searchQueryState={%22isMapVisible%22%3Atrue%2C%22mapBounds%22%3A{%22west%22%3A-105.91190088989258%2C%22east%22%3A-105.64513911010742%2C%22south%22%3A39.88949772255962%2C%22north%22%3A39.967295787905705}%2C%22filterState%22%3A{%22sort%22%3A{%22value%22%3A%22globalrelevanceex%22}}}

  const url = new URL("https://www.zillow.com/homes/for_sale/");
  url.searchParams.append(
    "searchQueryState",
    JSON.stringify({
      isMapVisible: true,
      mapBounds: {
        west: bounds.left(),
        east: bounds.right(),
        south: bounds.bottom(),
        north: bounds.top(),
      },
    }),
  );

  return url.toString();
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
  "#44AA99", // Teal
  "#117733", // Dark Green
  "#999933", // Olive Green
  "#AA4499", // Purple
  "#DDCC77", // Light Tan
  "#882255", // Dark Red
  "#332288", // Dark Blue
];

const keywords = [
  { query: "nwr[name=Costco]", radius: 5000 },
  { query: "nwr[amenity=cafe][name!=Starbucks]", radius: 1000 },
  { query: "nwr[amenity=school]", radius: 5000 },
];

assert.stab("start");

keywords.forEach((keyword) => {
  keyword.results = geo.query(`${keyword.query}(prefix=colorado)`);
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
      const color = colorPalette[index % colorPalette.length];

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
          "fill": colorPalette[index % colorPalette.length],
          "fill-opacity": 0.5,
          "url": zillowURL(bounds.asBound()),
        }),
      ],
    );
  }),
};

assert.stab("assert");
assert.geoJSON(payload);

assert.stab("return");
return payload;
