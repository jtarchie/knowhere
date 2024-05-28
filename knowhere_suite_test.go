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

var _ = Describe("Running the application", Ordered, func() {
	var path string

	cli := func(args ...string) *gexec.Session {
		command := exec.Command(path, args...)

		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())

		return session
	}

	BeforeAll(func() {
		var err error

		path, err = gexec.Build("github.com/jtarchie/knowhere", "--tags", "fts5", "-race")
		Expect(err).NotTo(HaveOccurred())
	})

	It("can generate a query", func() {
		session := cli(
			"query",
			`n[name="Hatfield Tunnel"](prefix="test")`,
		)

		Eventually(session, "5s").Should(gexec.Exit(0))
		Eventually(session.Out).Should(gbytes.Say("SELECT"))
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
	})

	When("serving the HTTP server", Ordered, func() {
		var (
			session *gexec.Session
			port    int
		)

		BeforeAll(func() {
			var err error
			port, err = freeport.GetFreePort()
			Expect(err).NotTo(HaveOccurred())

			buildPath, err := os.MkdirTemp("", "")
			Expect(err).NotTo(HaveOccurred())

			dbFilename := filepath.Join(buildPath, "test.sqlite")

			session = cli(
				"convert",
				"--osm", "./fixtures/sample.pbf",
				"--db", dbFilename,
				"--prefix", "test",
			)

			Eventually(session, "5s").Should(gexec.Exit(0))
			Expect(dbFilename).To(BeAnExistingFile())

			session = cli("server", "--port", strconv.Itoa(port), "--db", dbFilename)

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

		AfterAll(func() {
			session.Kill()
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
		})

		When("hitting the runtime endpoint", func() {
			It("returns the result in JSON", func() {
				client := req.C()
				response, err := client.R().
					SetRetryCount(3).
					SetBodyString(`
						const results = execute('nw[name="Hatfield Tunnel"](prefix="test")') ;
						return results.map((result) => result.name)
					`).
					Get(fmt.Sprintf("http://localhost:%d/api/runtime", port))

				Expect(err).NotTo(HaveOccurred())

				payload := &strings.Builder{}

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
				Entry("no source provided", ``, `source not provided in request body`),
				Entry("invalid javascript", `asdf;`, "evaluation error: ReferenceError: asdf is not defined at \u003ceval\u003e:3:5(1)"),
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
})
