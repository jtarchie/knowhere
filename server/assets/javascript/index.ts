document.addEventListener('DOMContentLoaded', function () {
  var searchInput = document.getElementById('search');

  // Function to handle the enter key press event
  function handleKeyPress(event) {
    if (event.key === 'Enter') {
      event.preventDefault(); // Prevents the default action of the enter key

      // Fetch data from API
      fetch('/api/search', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded',
        },
        body: 'search=' + encodeURIComponent(searchInput?.value)
      })
        .then(response => response.json())
        .then(data => {
          displayBoundingBoxes(data); // Log the response data
        })
        .catch(error => {
          console.error('Error', error);
        });
    }
  }

  // Add the event listener for the enter key press
  searchInput.addEventListener('keypress', handleKeyPress);
});


mapboxgl.accessToken = 'pk.eyJ1IjoianRhcmNoaSIsImEiOiJjbHBobmx0YWQwOG01MmlxeDAydGxlN2c5In0.o3yTh6k7uo_e3CBi_32R9Q';
var map = new mapboxgl.Map({
  container: 'map',
  style: 'mapbox://styles/mapbox/streets-v11', // Choose your style
  bounds: [
    [-125.0011, 24.9493], // Southwest coordinates (longitude, latitude)
    [-66.9326, 49.5904]   // Northeast coordinates
  ]
});

function displayBoundingBoxes(data) {
  var features = data.flatMap(box => {
    if (box.minLat === box.maxLat && box.minLon === box.maxLon) {
      return {
        type: 'Feature',
        geometry: {
          type: 'Point',
          coordinates: [box.minLon, box.minLat]
        }
      };
    } else {
      return [{
        type: 'Feature',
        geometry: {
          type: 'Point',
          coordinates: [box.minLon, box.minLat]
        }
      }, {
        type: 'Feature',
        geometry: {
          type: 'Polygon',
          coordinates: [
            [
              [box.minLon, box.minLat],
              [box.maxLon, box.minLat],
              [box.maxLon, box.maxLat],
              [box.minLon, box.maxLat],
              [box.minLon, box.minLat]
            ]
          ]
        }
      }];
    }
  });

  var source = map.getSource('results');
  if (source) {
    source.setData({
      type: 'FeatureCollection',
      features: features // Your new features
    });
  } else {
    map.addSource('results', {
      'type': 'geojson',
      'data': {
        'type': 'FeatureCollection',
        'features': features
      }
    });
  }

  map.addLayer({
    id: 'boundaries',
    type: 'fill',
    source: 'results',
    paint: {
      'fill-color': '#B42222',
    },
    'filter': ['==', '$type', 'Polygon']
  });

  map.addLayer({
    id: 'points',
    type: 'circle',
    source: 'results',
    paint: {
      'circle-radius': 6,
      'circle-color': '#B42222'
    },
    'filter': ['==', '$type', 'Point']
  });

  var bounds = new mapboxgl.LngLatBounds();
  features.forEach(feature => {
    if (feature.geometry.type === 'Point') {
      bounds.extend(feature.geometry.coordinates);
    } else if (feature.geometry.type === 'Polygon') {
      feature.geometry.coordinates[0].forEach(coord => {
        bounds.extend(coord);
      });
    }
  });

  if (features.length > 0) {
    map.fitBounds(bounds, {
      padding: 20
    });
  }
}