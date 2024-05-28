package main_test

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/imroc/req/v3"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"github.com/phayes/freeport"
)

func TestKnowhere(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Knowhere Suite")
}

var _ = Describe("Running the application", func() {
	var path string

	cli := func(args ...string) *gexec.Session {
		command := exec.Command(path, args...)

		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())

		return session
	}

	BeforeEach(func() {
		var err error

		path, err = gexec.Build("github.com/jtarchie/knowhere", "--tags", "fts5")
		Expect(err).NotTo(HaveOccurred())
	})

	It("can build sqlite file from osm pbf", func() {
		go func() {
			defer GinkgoRecover()

			http.Handle("/", http.FileServer(http.Dir("./fixtures")))
			err := http.ListenAndServe(":8848", nil)
			Expect(err).NotTo(HaveOccurred())
		}()

		buildPath, err := os.MkdirTemp("", "")
		Expect(err).NotTo(HaveOccurred())

		dbFilename := filepath.Join(buildPath, "test.sqlite")

		By("building all the things")

		session := cli(
			"build",
			"--config", "./fixtures/config.txt",
			"--db", dbFilename,
		)

		Eventually(session, "5s").Should(gexec.Exit(0))
		Expect(dbFilename).To(BeAnExistingFile())

		session = cli(
			"convert",
			"--osm", "./fixtures/sample.pbf",
			"--db", dbFilename,
			"--prefix", "test",
		)

		Eventually(session, "5s").Should(gexec.Exit(0))
		Expect(dbFilename).To(BeAnExistingFile())

		By("can generate a query")

		session = cli(
			"query",
			`n[name="Hatfield Tunnel"](prefix="test")`,
		)

		Eventually(session, "5s").Should(gexec.Exit(0))
		Eventually(session.Out).Should(gbytes.Say("SELECT"))

		By("serving HTTP")

		port, err := freeport.GetFreePort()
		Expect(err).NotTo(HaveOccurred())

		session = cli("server", "--port", strconv.Itoa(port), "--db", dbFilename)
		defer session.Kill()

		client := req.C()
		Eventually(func() error {
			_, err := client.R().
				SetRetryCount(3).
				Get(fmt.Sprintf("http://localhost:%d/", port))

			//nolint: wrapcheck
			return err
		}).ShouldNot(HaveOccurred())

		By("hitting search API endpoint")

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
						"id": 294,
						"type": "Feature",
						"geometry": {
							"type": "Point",
							"coordinates": [
                -0.24156,
                51.76005
							]
						},
            "properties": {
              "title": "Hatfield Tunnel"
            }
					}
				],
				"type": "FeatureCollection"
			}
		`))

		By("hitting the runtime API endpoint")

		response, err = client.R().
			SetRetryCount(3).
			SetBodyString(`
				const results = execute('nw[name="Hatfield Tunnel"](prefix="test")') ;
				return results.map((result) => result.name)
			`).
			Get(fmt.Sprintf("http://localhost:%d/api/runtime", port))

		Expect(err).NotTo(HaveOccurred())

		payload = &strings.Builder{}

		_, err = io.Copy(payload, response.Body)
		Expect(err).NotTo(HaveOccurred())

		Expect(payload.String()).To(MatchJSON(`["Hatfield Tunnel"]`))

		By("hitting prefixes API endpoint")

		response, err = client.R().
			SetRetryCount(3).
			Get(fmt.Sprintf("http://localhost:%d/api/prefixes", port))

		Expect(err).NotTo(HaveOccurred())

		payload = &strings.Builder{}

		_, err = io.Copy(payload, response.Body)
		Expect(err).NotTo(HaveOccurred())

		Expect(payload.String()).To(MatchJSON(`
		{
			"prefixes": [
				{
					"name": "Sample Some Long Name",
					"slug": "sample_some_long_name",
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
				},
				{
					"name": "Sample",
					"slug": "sample",
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
				},
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
