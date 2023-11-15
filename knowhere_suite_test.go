package main_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

func TestKnowhere(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Knowhere Suite")
}

func cli(args ...string) string {
	path, err := gexec.Build("github.com/jtarchie/knowhere", "--tags", "fts5")
	Expect(err).NotTo(HaveOccurred())

	command := exec.Command(path, args...)

	session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
	Expect(err).ToNot(HaveOccurred())

	Eventually(session, "5s").Should(gexec.Exit(0))

	return string(session.Out.Contents())
}

var _ = Describe("Running the application", func() {
	It("can build sqlite file from osm pbf", func() {
		buildPath, err := os.MkdirTemp("", "")
		Expect(err).NotTo(HaveOccurred())

		dbFilename := filepath.Join(buildPath, "test.sqlite")

		_ = cli(
			"build",
			"--osm", "./fixtures/sample.pbf",
			"--db", dbFilename,
		)

		By("can generate a query")

		output := cli(
			"query",
			"n[name=Starbucks]",
		)

		Expect(output).To(ContainSubstring("SELECT"))
	})
})
