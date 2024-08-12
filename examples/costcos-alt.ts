/// <reference path="../docs/src/global.d.ts" />

const allCostcos = query.execute(`nwr[name=~Costco](prefix=colorado)`);

const entries = allCostcos.cluster(2000);

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
