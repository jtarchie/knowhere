/// <reference path="../docs/src/global.d.ts" />

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

const keywords = [
  {
    radius: 5000,
    results: query.execute(`nwr[name=~Costco](area=colorado)`),
  },
  {
    radius: 1000,
    results: query.execute(
      `nwr[amenity=cafe][name][name!~Starbucks](area=colorado)`,
    ),
  },
  {
    radius: 5000,
    results: query.execute(`nwr[amenity=school][name](area=colorado)`),
  },
];

keywords.sort((a, b) => a.results.length - b.results.length);

const neighbors = keywords[0].results.cluster(500).map((
  entry,
) => geo.asResults(entry));

keywords.slice(1).forEach((keyword) => {
  const tree = keyword.results.asTree(keyword.radius);

  neighbors.forEach((entries) => {
    const extended = entries[0].bound().extend(keywords[0].radius);

    const nearby = tree.nearby(extended, 1);
    if (nearby.length === 1) {
      entries.push(nearby[0]);
    }
  });
});

const payload = {
  type: "FeatureCollection",
  features: neighbors.flatMap((entries, index) => {
    if (entries.length !== keywords.length) {
      return;
    }

    const features = entries.flatMap((entry, index) => {
      const color = colors.pick(index);

      const feature = entry.asFeature({
        "marker-color": color,
        index: index,
      });

      return feature;
    });

    const bounds = geo.asBounds(
      ...entries.map((entry, index) =>
        entry.bound().extend(keywords[index].radius)
      ),
    );

    return features.concat(
      [
        bounds.asFeature({
          "fill": colors.pick(index),
          "fill-opacity": 0.5,
          "url": zillowURL(bounds.asBound()),
        }),
      ],
    );
  }).filter(Boolean),
};

assert.geoJSON(payload);

export { payload };
