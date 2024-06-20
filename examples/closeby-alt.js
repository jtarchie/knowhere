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

assert.stab("start");

const keywords = [
  { query: "nwr[name=Costco]", radius: 5000 },
  { query: "nwr[amenity=cafe][name][name!=Starbucks]", radius: 1000 },
  { query: "nwr[amenity=school][name]", radius: 5000 },
];

keywords.forEach((keyword) => {
  keyword.results = geo.query(`${keyword.query}(prefix=colorado)`);
});

keywords.sort((a, b) => a.results.length - b.results.length);

assert.stab("query");

const neighbors = new Map();
const cluster = keywords[0].results.cluster(500);
cluster.forEach((entry) => {
  neighbors.set(entry.id, new Map());
});

assert.stab("cluster");

const expectedNeighbors = 2;

keywords.slice(1).forEach((keyword) => {
  const grouped = cluster.overlap(
    keyword.results,
    keywords[0].radius,
    expectedNeighbors - 1,
  );
  grouped.forEach((values) => {
    values.forEach((value) => neighbors.get(values[0].id).set(value.id, value));
  });
});

assert.stab("nearby");

const payload = {
  type: "FeatureCollection",
  features: [...neighbors.values()].flatMap((set, index) => {
    const entries = [...set.values()];

    if (entries.length !== keywords.length) {
      return;
    }

    const features = entries.flatMap((entry, index) => {
      const color = geo.color(index);

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
          "fill": geo.color(index),
          "fill-opacity": 0.5,
          "url": zillowURL(bounds.asBound()),
        }),
      ],
    );
  }).filter(Boolean),
};

assert.stab("payload");

assert.geoJSON(payload);

assert.stab("assert GeoJSON");

return payload;
