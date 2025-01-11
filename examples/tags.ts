/// <reference path="../docs/src/global.d.ts" />

const areas = query.execute(
  `nwr[boundary=administrative][admin_level>=6][name=~Denver](area="colorado")`,
);
assert.eq(areas.length == 1, "one area expected");

const area = areas[0];
const bounds = area.bound().extend(20_000);
const entries = query.execute(
  `nwr[name](area=colorado)(bb=${bounds.asBB()})`, // bb=minLon,minLat,maxLon,maxLat
);
assert.eq(entries.length > 0, "entries expected");

const tagCountsMap = entries.tagCount();
const tagCounts = Object.entries(tagCountsMap).sort(([, a], [, b]) => b - a);

const excludePattern =
  /geohash|name|tiger|source|fixme|_id|attribution|addr|wikipedia|url|gtfs/i;
const payload = tagCounts.slice(0, 500).filter(([tag]) =>
  !excludePattern.test(tag)
).slice(0, 100).map(([tag]) => tag);

export { payload };
