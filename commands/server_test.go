package commands_test

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/imroc/req/v3"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/phayes/freeport"

	"github.com/jtarchie/knowhere/commands"
)

var _ = Describe("Server", func() {
	var (
		port int
	)

	BeforeEach(func() {
		var err error
		port, err = freeport.GetFreePort()
		Expect(err).NotTo(HaveOccurred())

		buildPath, err := os.MkdirTemp("", "")
		Expect(err).NotTo(HaveOccurred())

		dbFilename := filepath.Join(buildPath, "test.sqlite")

		convert := &commands.Convert{
			OSM:    "../fixtures/sample.pbf",
			DB:     dbFilename,
			Prefix: "test",
		}
		err = convert.Run(GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		Expect(dbFilename).To(BeAnExistingFile())

		server := &commands.Server{
			Port:           port,
			DB:             dbFilename,
			RuntimeTimeout: time.Second,
		}
		go func() {
			defer GinkgoRecover()

			err := server.Run(GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
		}()

		// wait for it to start
		client := req.C()
		Eventually(func() error {
			_, err := client.R().
				SetRetryCount(3).
				Get(fmt.Sprintf("http://localhost:%d/", port))

			//nolint: wrapcheck
			return err
		}).ShouldNot(HaveOccurred())
	})

	It("hitting search API endpoint", func() {
		client := req.C()
		response, err := client.R().
			SetRetryCount(3).
			AddQueryParam("search", `nw[name="Hatfield Tunnel"](prefix="test")`).
			Get(fmt.Sprintf("http://localhost:%d/api/search", port))

		Expect(err).NotTo(HaveOccurred())

		payload := &strings.Builder{}

		_, err = io.Copy(payload, response.Body)
		Expect(err).NotTo(HaveOccurred())

		Expect(payload.String()).To(MatchJSON(`
		{
        "features": [
          {
            "type": "Feature",
            "geometry": {
              "type": "Point",
              "coordinates": [
                -0.23752,
                51.76459
              ]
            },
            "properties": {
              "id": 294,
              "title": "Hatfield Tunnel",
              "type": 2
            }
          }
        ],
        "type": "FeatureCollection"
      }
	`))
	})

	When("hitting the runtime endpoint", func() {
		It("returns the result in JSON", func() {
			source := `
			const results = geo.query('nw[name="Hatfield Tunnel"](prefix="test")') ;
			return results.map((result) => result.name)
			`
			client := req.C()
			response, err := client.R().
				SetRetryCount(3).
				SetBodyString(source).
				Get(fmt.Sprintf("http://localhost:%d/api/runtime", port))

			Expect(err).NotTo(HaveOccurred())

			payload := &strings.Builder{}

			_, err = io.Copy(payload, response.Body)
			Expect(err).NotTo(HaveOccurred())

			Expect(payload.String()).To(MatchJSON(`["Hatfield Tunnel"]`))

			response, err = client.R().
				SetRetryCount(3).
				SetQueryParam("source", source).
				Get(fmt.Sprintf("http://localhost:%d/api/runtime", port))

			Expect(err).NotTo(HaveOccurred())

			payload = &strings.Builder{}

			_, err = io.Copy(payload, response.Body)
			Expect(err).NotTo(HaveOccurred())

			Expect(payload.String()).To(MatchJSON(`["Hatfield Tunnel"]`))
		})

		DescribeTable("returns error payload on exception", func(source string, errMsg string) {
			client := req.C()
			response, err := client.R().
				SetRetryCount(3).
				SetBodyString(source).
				Get(fmt.Sprintf("http://localhost:%d/api/runtime", port))

			Expect(err).NotTo(HaveOccurred())

			payload := &strings.Builder{}

			_, err = io.Copy(payload, response.Body)
			Expect(err).NotTo(HaveOccurred())

			Expect(payload.String()).To(MatchJSON(fmt.Sprintf(`{"error":%q}`, errMsg)))
		},
			Entry("timeouts on infinite loop", "for(;;) {}", `evaluation error: vm timed out at main.js:1:3(1)`),
			Entry("no source provided", ``, `source not provided in request body`),
			Entry("invalid javascript", `asdf;`, "evaluation error: ReferenceError: asdf is not defined at main.js:1:15(1)"),
			Entry("assertion fail", `assert.eq(false, "this did not work")`, "evaluation error: assertion failed: this did not work at main.js:1:24(6)"),
			Entry("assertion geojson", `assert.geoJSON({})`, "evaluation error: assert of geojson failed: missing type at main.js:1:29(5)"),
		)
	})

	It("hitting prefixes API endpoint", func() {
		client := req.C()
		response, err := client.R().
			SetRetryCount(3).
			Get(fmt.Sprintf("http://localhost:%d/api/prefixes", port))

		Expect(err).NotTo(HaveOccurred())

		payload := &strings.Builder{}

		_, err = io.Copy(payload, response.Body)
		Expect(err).NotTo(HaveOccurred())

		Expect(payload.String()).To(MatchJSON(`
		{
			"prefixes": [
				{
					"name": "Test",
					"slug": "test",
					"bounds": [
						[
							-0.24156,
							51.76005
						],
						[
							-0.21629,
							51.77425
						]
					]
				}
			]
		}
	`))
	})
})
