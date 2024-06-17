const prefixes = geo.prefixes();

const allColleges = geo.asResults(
  ...prefixes.flatMap((prefix) => {
    return geo.query(
      `wr[amenity=university,college][name=university,college](prefix=${prefix.name})`,
    );
  }),
);

const radius = 500
const overlap = 2000

const clustered = allColleges.cluster(radius);
const grouped = clustered.overlap(clustered, overlap, 4);

const payload = {
  type: "FeatureCollection",
  features: grouped.flatMap((entries, index) => {
    const features = entries.flatMap((entry) => {
      const feature = entry.asFeature({
        "marker-color": geo.color(index),
        index: index,
      });

      return feature;
    });

    const bounds = geo.asBounds(
      ...entries.map((entry, index) => entry.bbox().extend(overlap)),
    );

    return features.concat(
      [
        bounds.asFeature({
          "fill": geo.color(index),
          "fill-opacity": 0.2,
        }),
      ],
    );
  }),
};

return payload;
