const prefixes = query.prefixes();

const allUnis = geo.asResults(
  ...prefixes.flatMap((prefix) => {
    return query.execute(
      `wr[amenity=university][name](prefix=${prefix.name})`,
    );
  }),
);

const radius = 500;
const overlap = 3000;

const clustered = allUnis.cluster(radius);
const grouped = clustered.overlap(clustered, overlap, 0, 3);

const payload = {
  type: "FeatureCollection",
  features: grouped.flatMap((entries, index) => {
    const features = entries.flatMap((entry) => {
      const feature = entry.asFeature({
        "marker-color": colors.pick(index),
        index: index,
      });

      return feature;
    });

    const bounds = geo.asBounds(
      ...entries.map((entry) => entry.bbox().extend(overlap)),
    );

    return features.concat(
      [
        bounds.asFeature({
          "fill": colors.pick(index),
          "fill-opacity": 0.2,
        }),
      ],
    );
  }),
};

return payload;
