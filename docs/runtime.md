# Embedded JavaScript Runtime Documentation

This document provides an overview of the available functions in the embedded
JavaScript runtime for a Golang program. The functions are inferred from the
given usage examples and are documented with simple examples to demonstrate
their usage.

Every script must `return` the final value, which will be marshalled to JSON. A
script can be invoked via CLI `knowhere runtime` or the endpoint `/api/runtime`.
See `examples/` for more.

## Functions

### `query.execute(query)`

Performs a geographical query.

**Parameters:**

- `query`: A string representing the query.

**Returns:**

- An array of results matching the query.

**Example:**

```javascript
const results = query.execute("nwr[name=~Costco](prefix=colorado)");
console.log(results); // Outputs results for Costcos in Colorado
```

### `query.prefixes()`

Retrieves available geographical prefixes.

**Returns:**

- An array of prefix objects.

**Example:**

```javascript
const prefixes = query.prefixes();
console.log(prefixes); // Outputs available geographical prefixes
```

### `geo.asResults(...queries)`

Combines multiple queries into a single result set.

**Parameters:**

- `queries`: Multiple query results.

**Returns:**

- A combined result set.

**Example:**

```javascript
const allUnis = geo.asResults(
  ...prefixes.flatMap((prefix) => {
    return query.execute(`wr[amenity=university][name](prefix=${prefix.name})`);
  }),
);
console.log(allUnis); // Outputs combined results for universities
```

### `colors.pick(index)`

Generates a color based on an index.

**Parameters:**

- `index`: A numerical index.

**Returns:**

- A color string.

**Example:**

```javascript
const color = colors.pick(1);
console.log(color); // Outputs a color string based on the index
```

### `geo.asBounds(...entries)`

Creates a bounding box from multiple entries.

**Parameters:**

- `entries`: Multiple entries to be included in the bounding box.

**Returns:**

- A bounding box object.

**Example:**

```javascript
const bounds = geo.asBounds(entry1, entry2, entry3);
console.log(bounds); // Outputs a bounding box object
```

### `assert.stab(message)`

Inserts a stable checkpoint for debugging purposes.

**Parameters:**

- `message`: A string message indicating the checkpoint.

**Example:**

```javascript
assert.stab("start");
// Some code here
assert.stab("query");
```

### `assert.eq(value1, value2, message)`

Asserts that two values are equal.

**Parameters:**

- `value1`: The first value.
- `value2`: The second value.
- `message`: A string message indicating the assertion.

**Example:**

```javascript
assert.eq(5, 5, "expected 5 to equal 5");
```

### `assert.geoJSON(payload)`

Asserts that a payload is valid GeoJSON.

**Parameters:**

- `payload`: The GeoJSON payload.

**Example:**

```javascript
const payload = {
  type: "FeatureCollection",
  features: [],
};

assert.geoJSON(payload); // Asserts the payload is valid GeoJSON
```

This document provides an overview of the available functions and how to use
them in the embedded JavaScript runtime within a Golang program. The examples
illustrate practical applications of these functions for geographical queries
and data manipulation.
