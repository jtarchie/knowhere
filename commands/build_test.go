package commands_test

import (
	"net/http"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/jtarchie/knowhere/commands"
)

var _ = Describe("Build", func() {
	It("can build sqlite file from osm pbf", func() {
		go func() {
			defer GinkgoRecover()

			http.Handle("/", http.FileServer(http.Dir("../fixtures")))
			err := http.ListenAndServe(":8848", nil)
			Expect(err).NotTo(HaveOccurred())
		}()

		buildPath, err := os.MkdirTemp("", "")
		Expect(err).NotTo(HaveOccurred())

		dbFilename := filepath.Join(buildPath, "test.sqlite")

		By("building all the things")

		build := &commands.Build{
			Config: "../fixtures/config.txt",
			DB:     dbFilename,
		}
		err = build.Run(GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		Expect(dbFilename).To(BeAnExistingFile())

		convert := &commands.Convert{
			OSM:    "../fixtures/sample.pbf",
			DB:     dbFilename,
			Prefix: "test",
		}
		err = convert.Run(GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		Expect(dbFilename).To(BeAnExistingFile())
	})
})
