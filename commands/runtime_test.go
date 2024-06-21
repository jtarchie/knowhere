package commands_test

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/jtarchie/knowhere/commands"
)

var _ = Describe("Runtime", func() {
	When("using the examples", func() {
		dbFilename := "../.build/entries.db"

		BeforeEach(func() {
			if _, err := os.Stat(dbFilename); errors.Is(err, os.ErrNotExist) {
				Skip(".build/entries.db does not exist")
			}
		})

		It("should have some examples", func() {
			matches, err := filepath.Glob("../examples/*.js")
			Expect(err).NotTo(HaveOccurred())
			Expect(len(matches)).NotTo(Equal(0))
		})

		matches, _ := filepath.Glob("../examples/*.js")
		for _, match := range matches {
			It(fmt.Sprintf("ensures that %q passes", match), func() {
				file, err := os.Open(match)
				Expect(err).NotTo(HaveOccurred())

				runtime := &commands.Runtime{
					Filename:       file,
					DB:             dbFilename,
					RuntimeTimeout: 10 * time.Second,
				}

				err = runtime.Run(GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())
			})
		}
	})
})
