const allCostcos = execute(`nwr[name=Costco](prefix=colorado)`);

const bounds = [];

const entries = allCostcos.filter((costco) => {
  const extended = costco.bbox().extend(2);

  if (bounds.some((bbox) => bbox.intersects(extended))) {
    return false;
  }

  bounds.push(extended);
  return true;
});

assert.eq(allCostcos.length > entries.length, "expected fewer entries");

const payload = {
  type: "FeatureCollection",
  features: entries.map((entry) => {
    return entry.asFeature();
  }),
};

assert.geoJSON(payload);

return payload;
