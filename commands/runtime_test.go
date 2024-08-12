package commands_test

import (
	"bytes"
	"encoding/json"
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
			matches, err := filepath.Glob("../examples/*.ts")
			Expect(err).NotTo(HaveOccurred())
			Expect(len(matches)).NotTo(Equal(0))
		})

		examples, _ := filepath.Glob("../examples/*.ts")
		docs, _ := filepath.Glob("../docs/src/examples/*.ts")

		for _, match := range append(examples, docs...) {
			It(fmt.Sprintf("ensures that %q passes", match), func() {
				file, err := os.Open(match)
				Expect(err).NotTo(HaveOccurred())

				runtime := &commands.Runtime{
					Filename:       file,
					DB:             dbFilename,
					RuntimeTimeout: 10 * time.Second,
				}

				buffer := &bytes.Buffer{}
				err = runtime.Run(buffer)
				Expect(err).NotTo(HaveOccurred())

				var payload interface{}
				err = json.Unmarshal(buffer.Bytes(), &payload)
				Expect(payload).NotTo(BeNil())
				Expect(err).NotTo(HaveOccurred())
			})
		}
	})
})
