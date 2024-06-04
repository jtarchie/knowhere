package main_test

import (
	"errors"
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
				Entry("no source provided", ``, `source not provided in request body`),
				Entry("invalid javascript", `asdf;`, "evaluation error: ReferenceError: asdf is not defined at main.js:1:15(1)"),
				Entry("assertion fail", `assert.eq(false, "this did not work")`, "evaluation error: assertion failed: this did not work at main.js:1:24(6)"),
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

	When("using the examples", func() {
		var (
			session *gexec.Session
			port    int
		)

		BeforeEach(func() {
			var err error
			port, err = freeport.GetFreePort()
			Expect(err).NotTo(HaveOccurred())

			dbFilename := "./.build/entries.db"
			if _, err := os.Stat(dbFilename); errors.Is(err, os.ErrNotExist) {
				Skip(".build/entries.db does not exist")
			}

			session = cli("server", "--port", strconv.Itoa(port), "--db", dbFilename)
		})

		AfterEach(func() {
			session.Kill()
		})

		It("ensures that they all pass", func() {
			matches, err := filepath.Glob("./examples/*.js")
			Expect(err).NotTo(HaveOccurred())
			Expect(len(matches)).To(BeNumerically(">", 0))

			for _, match := range matches {
				contents, err := os.ReadFile(match)
				Expect(err).NotTo(HaveOccurred())

				// result := api.Transform(string(contents), api.TransformOptions{
				// 	MinifyWhitespace:  true,
				// 	MinifyIdentifiers: true,
				// 	MinifySyntax:      true,
				// 	Sourcemap:         api.SourceMapInline,
				// })
				// Expect(result.Errors).To(HaveLen(0))

				client := req.C()
				response, err := client.R().
					SetRetryCount(3).
					SetBodyString(string(contents)).
					Get(fmt.Sprintf("http://localhost:%d/api/runtime", port))

				Expect(err).NotTo(HaveOccurred())
				Expect(response.StatusCode).To(Equal(http.StatusOK))
			}
		})
	})
})
