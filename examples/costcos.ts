/// <reference path="../docs/examples/global.d.ts" />

const allCostcos = query.execute(`nwr[name=~Costco](prefix=colorado)`);

const bounds: Bound[] = [];

const entries = allCostcos.filter((costco) => {
  const extended = costco.bound().extend(2000);

  if (bounds.some((bbox) => bbox.intersects(extended))) {
    return false;
  }

  bounds.push(extended);
  return true;
});

assert.eq(allCostcos.length > entries.length, "expected fewer entries");

const payload = {
  type: "FeatureCollection",
  features: entries.map((entry, index) => {
    return entry.asFeature({
      "marker-color": colors.pick(index),
    });
  }),
};

assert.geoJSON(payload);

export { payload };
