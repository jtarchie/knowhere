package services_test

import (
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/jtarchie/knowhere/services"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("When using the runtime", func() {
	var client *sql.DB

	BeforeEach(func() {
		tmpPath, err := os.MkdirTemp("", "")
		Expect(err).NotTo(HaveOccurred())

		dbPath := filepath.Join(tmpPath, "test.db")

		converter := services.NewConverter(
			"../fixtures/sample.pbf",
			dbPath,
			"test",
			[]string{"*"},
		)
		err = converter.Execute()
		Expect(err).NotTo(HaveOccurred())

		client, err = sql.Open("sqlite3", dbPath)
		Expect(err).NotTo(HaveOccurred())
	})

	It("can run hello world", func() {
		runtime := services.NewRuntime(client, time.Second)
		value, err := runtime.Execute(`
			return "Hello, World"
		`)
		Expect(err).NotTo(HaveOccurred())
		Expect(value.Export()).To(BeEquivalentTo("Hello, World"))
	})

	When("using the bounding box", func() {
		It("returns the original", func() {
			runtime := services.NewRuntime(client, time.Second)
			value, err := runtime.Execute(`
				const results = geo.query('nw[name="Hatfield Tunnel"](prefix=test)');
				assert.eq(results.length == 1);
				
				return results[0].bbox()
			`)
			Expect(err).NotTo(HaveOccurred())
			Expect(toJSON(value.Export())).To(MatchJSON(`{
        "Min": [
          -0.24156,
          51.76005
        ],
        "Max": [
          -0.23348,
          51.76913
        ]
      }`))
		})
	})

	It("asserts valid GeoJSON", func() {
		runtime := services.NewRuntime(client, time.Second)
		_, err := runtime.Execute(`
			const payload = {};
			assert.geoJSON(payload);
			return payload;
		`)
		Expect(err).To(HaveOccurred())

		value, err := runtime.Execute(`
			const payload = {
				type: "Feature",
				geometry: {
					type: "Point",
					coordinates: [125.6, 10.1]
				},
				properties: {
					name: "Dinagat Islands"
				}
			};
			assert.geoJSON(payload);
			return payload;
		`)
		Expect(err).NotTo(HaveOccurred())

		contents := toJSON(value.Export())
		Expect(contents).To(MatchJSON(`
			{
				"type": "Feature",
				"geometry": {
					"type": "Point",
					"coordinates": [
						125.6,
						10.1
					]
				},
				"properties": {
					"name": "Dinagat Islands"
				}
			}`),
		)
	})
})

func toJSON(object any) string {
	contents, err := json.Marshal(object)
	Expect(err).NotTo(HaveOccurred())

	return string(contents)
}
