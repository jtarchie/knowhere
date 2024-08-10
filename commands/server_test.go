package commands_test

import (
	"encoding/json"
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
			CacheSize:      5000,
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
			const results = query.execute('nw[name="Hatfield Tunnel"](prefix="test")') ;
			assert.stab(JSON.stringify(results))
			return results.map((result) => result.tags['name'])
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

			var payloadMap map[string]string
			err = json.Unmarshal([]byte(payload.String()), &payloadMap)
			Expect(err).NotTo(HaveOccurred())
			Expect(payloadMap["error"]).To(ContainSubstring(errMsg))
		},
			Entry("timeouts on infinite loop", "for(;;) {}", `evaluation error: vm timed out at main.js`),
			Entry("no source provided", ``, `source not provided in request body`),
			Entry("invalid javascript", `asdf;`, "evaluation error: ReferenceError: asdf is not defined at main.js"),
			Entry("syntax error", `const a = a, b => {}`, "evaluation error: SyntaxError:"),
			Entry("assertion fail", `assert.eq(false, "this did not work")`, "evaluation error: assertion failed: this did not work at main.js"),
			Entry("assertion geojson", `assert.geoJSON({})`, "evaluation error: assert of geojson failed: missing type at main.js"),
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
