mapboxgl.accessToken =
  "pk.eyJ1IjoianRhcmNoaSIsImEiOiJjbHBobmx0YWQwOG01MmlxeDAydGxlN2c5In0.o3yTh6k7uo_e3CBi_32R9Q";
const map = new mapboxgl.Map({
  container: "map",
  style: "mapbox://styles/mapbox/streets-v12",
  bounds: [
    [-124.7844079, 24.396308],
    [-66.9513812, 49.384358],
  ],
});

map.addControl(
  new mapboxgl.NavigationControl(),
  "top-right",
);

const source = `
const entries = geo.query("nwr[name=Costco](prefix=colorado)");

const payload = {
  type: "FeatureCollection",
  features: entries.map((entry) => {
    return entry.asFeature();
  }),
};

return payload
`;

map.once("idle", () => {
  var bounds = new mapboxgl.LngLatBounds();

  map.querySourceFeatures("map-data").forEach(function (feature) {
    bounds.extend(feature.geometry.coordinates);
  });

  map.fitBounds(bounds);
});

map.on("load", () => {
  map.addSource("map-data", {
    type: "geojson",
    data: `/api/runtime?source=${encodeURIComponent(source)}`,
  });

  const color = "#56B4E9";

  map.addLayer({
    id: "map-data-marker",
    type: "circle",
    source: "map-data",
    paint: {
      "circle-radius": [
        "match",
        ["get", "marker-size"],
        "small",
        4,
        "medium",
        8,
        "large",
        12,
        8,
      ],
      "circle-color": ["coalesce", ["get", "marker-color"], color],
      "circle-stroke-width": 2,
      "circle-stroke-color": "#ffffff",
    },
    filter: ["==", ["geometry-type"], "Point"],
  });

  map.addLayer({
    id: "map-data-fill",
    type: "fill",
    source: "map-data",
    paint: {
      "fill-color": ["coalesce", ["get", "fill"], color],
      "fill-opacity": ["coalesce", ["get", "fill-opacity"], 0.3],
    },
    filter: ["==", ["geometry-type"], "Polygon"],
  });

  map.addLayer({
    id: "map-data-fill-outline",
    type: "line",
    source: "map-data",
    paint: {
      "line-color": ["coalesce", ["get", "stroke"], color],
      "line-width": ["coalesce", ["get", "stroke-width"], 2],
      "line-opacity": ["coalesce", ["get", "stroke-opacity"], 1],
    },
    filter: ["==", ["geometry-type"], "Polygon"],
  });

  map.addLayer({
    id: "map-data-fill",
    type: "fill",
    source: "map-data",
    paint: {
      "fill": ["coalesce", ["get", "fill"], color],
      "line-opacity": ["coalesce", ["get", "fill-opacity"], 1],
    },
    filter: ["==", ["geometry-type"], "Polygon"],
  });

  map.addLayer({
    id: "map-data-line",
    type: "line",
    source: "map-data",
    paint: {
      "line-color": ["coalesce", ["get", "stroke"], color],
      "line-width": ["coalesce", ["get", "stroke-width"], 2],
      "line-opacity": ["coalesce", ["get", "stroke-opacity"], 1],
    },
    filter: ["==", ["geometry-type"], "LineString"],
  });
});
