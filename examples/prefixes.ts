/// <reference path="../docs/global.d.ts" />

const prefixes = query.prefixes();
const prefix = prefixes[Math.floor(Math.random() * prefixes.length)];

const allCostcos = query.execute(`nwr[name=~Costco](prefix=${prefix.name})`);

const bounds: Bound[] = [];

const entries = allCostcos.filter((costco) => {
  const extended = costco.bound().extend(2000);

  if (bounds.some((bbox) => bbox.intersects(extended))) {
    return false;
  }

  bounds.push(extended);
  return true;
});

assert.eq(
  allCostcos.length >= entries.length,
  `expected ${allCostcos.length} >= ${entries.length}`,
);

const payload = {
  type: "FeatureCollection",
  features: entries.map((entry) => {
    return entry.asFeature();
  }),
};

assert.geoJSON(payload);

export { payload };
