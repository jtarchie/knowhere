package services_test

import (
	"encoding/json"

	"github.com/jtarchie/knowhere/services"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("When using the runtime", func() {
	It("can run hello world", func() {
		runtime := services.NewRuntime(nil)
		value, err := runtime.Execute(`
			return "Hello, World"
		`)
		Expect(err).NotTo(HaveOccurred())
		Expect(value).To(BeEquivalentTo("Hello, World"))
	})

	It("asserts valid GeoJSON", func() {
		runtime := services.NewRuntime(nil)
		_, err := runtime.Execute(`
			const payload = {};
			if (assertGeoJSON(payload) === false) {
				throw "bork"
			}
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
			if (assertGeoJSON(payload) === false) {
				throw "bork";
			}
			return payload;
		`)
		Expect(err).NotTo(HaveOccurred())

		contents, err := json.Marshal(value)
		Expect(err).NotTo(HaveOccurred())

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
