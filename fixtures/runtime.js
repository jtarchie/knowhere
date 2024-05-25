const starbucks = execute(`n[name=Starbucks][prefix=colorado]`);
const coffees = execute(`n[name=coffee][prefix=colorado]`);

const results = starbucks.slice(0, 2).concat(coffees.slice(0, 2));

return {
  type: "FeatureCollection",
  features: results.map((result) => {
    const payload = turf.bboxPolygon(turf.bbox(turf.point([result.MinLon, result.MinLat])))
    return {
      ...payload, ...{
        id: result.ID,
        type: "Feature",
        properties: {
          name: result.Name,
        }
      }
    }
  })
}
