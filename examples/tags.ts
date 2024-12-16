/// <reference path="../docs/src/global.d.ts" />

const areas = query.execute(
  `nwr[boundary=administrative][admin_level>=6][name=~Denver](area="colorado")`,
);
assert.eq(areas.length == 1, "one area expected");

const area = areas[0];
const bounds = area.bound().extend(20_000);
// min and max are [lon, lat]
const min = bounds.min();
const max = bounds.max();
// bb=minLon,minLat,maxLon,maxLat
const entries = query.execute(
  `nwr[name](area=colorado)(bb=${min[0]},${min[1]},${max[0]},${max[1]})`,
);
assert.eq(entries.length > 0, "entries expected");

const tagCounts = Object.entries(
  entries.flatMap((entry) => {
    const tags = Object.keys(entry.tags);
    return tags;
  }).reduce(
    (
      acc: { [key: string]: number },
      item: string,
    ): { [key: string]: number } => {
      acc[item] = (acc[item] || 0) + 1;
      return acc;
    },
    {},
  ),
).sort((a, b) => b[1] - a[1]).filter((item) =>
  !/geohash|name|tiger|source|fixme|_id|attribution|addr|wikipedia|url|gtfs/i.test(item[0])
);

const payload = tagCounts.slice(0, 200);

export { payload };
