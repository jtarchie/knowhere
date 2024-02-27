package main_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
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
		buildPath, err := os.MkdirTemp("", "")
		Expect(err).NotTo(HaveOccurred())

		dbFilename := filepath.Join(buildPath, "test.sqlite")

		session := cli(
			"build",
			"--osm", "./fixtures/sample.pbf",
			"--db", dbFilename,
			"--prefix", "test",
		)

		Eventually(session, "5s").Should(gexec.Exit(0))
		Expect(dbFilename).To(BeAnExistingFile())

		By("can generate a query")

		session = cli(
			"query",
			"n[name=Starbucks]",
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
	})
})
