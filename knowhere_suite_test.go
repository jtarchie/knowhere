package main_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/alecthomas/kong"
	"github.com/jtarchie/knowhere/commands"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestKnowhere(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Knowhere Suite")
}

func cli(args ...string) {
	cli := &commands.CLI{}
	parser, err := kong.New(cli)
	Expect(err).NotTo(HaveOccurred())

	ctx, err := parser.Parse(args)
	Expect(err).NotTo(HaveOccurred())

	err = ctx.Run()
	Expect(err).NotTo(HaveOccurred())
}

var _ = Describe("Running the application", func() {
	It("can build sqlite file from osm pbf", func() {
		buildPath, err := os.MkdirTemp("", "")
		Expect(err).NotTo(HaveOccurred())

		dbFilename := filepath.Join(buildPath, "test.sqlite")

		cli(
			"build",
			"--osm", "./fixtures/sample.pbf",
			"--db", dbFilename,
		)

		Expect(dbFilename).To(BeAnExistingFile())
	})
})
