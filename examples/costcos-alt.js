assert.stab("start");

const allCostcos = geo.query(`nwr[name=~Costco](prefix=colorado)`);

assert.stab("query");

const entries = allCostcos.cluster(2000);

assert.stab("cluster");

assert.eq(allCostcos.length > entries.length, "expected fewer entries");

const payload = {
  type: "FeatureCollection",
  features: entries.map((entry, index) => {
    return entry.asFeature({
      "marker-color": geo.color(index),
    });
  }),
};

assert.stab("payload");

assert.geoJSON(payload);

assert.stab("assert");
return payload;
