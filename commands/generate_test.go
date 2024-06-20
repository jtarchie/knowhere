package commands_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"

	"github.com/jtarchie/knowhere/commands"
)

var _ = Describe("Generate", func() {
	It("can generate a query", func() {
		generate := commands.Generate{
			Value: `n[name="Hatfield Tunnel"](prefix="test")`,
		}
		buffer := gbytes.NewBuffer()

		err := generate.Run(buffer)
		Expect(err).NotTo(HaveOccurred())
		Expect(buffer).Should(gbytes.Say("SELECT"))
	})
})
