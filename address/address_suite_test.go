package address_test

import (
	"os"
	"strings"
	"testing"

	"github.com/jtarchie/knowhere/address"
	"github.com/recursionpharma/go-csv-map"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestAddress(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Address Suite")
}

var _ = Describe("Parse", func() {
	It("parses US addresses to a threshold", func() {
		file, err := os.Open("./test_data.csv")
		Expect(err).NotTo(HaveOccurred())

		reader := csvmap.NewReader(file)
		reader.Columns, err = reader.ReadHeader()
		Expect(err).NotTo(HaveOccurred())

		records, err := reader.ReadAll()
		Expect(err).NotTo(HaveOccurred())

		valid := 0
		total := 0
		for _, record := range records {
			if record["country_code"] == "us" {
				total++
				fullAddress := record["full_address"]
				parsedAddress, ok := address.Parse(fullAddress)
				if ok {
					parseValid := true
					for key, value := range parsedAddress {
						if strings.ToLower(value) != record[key] {
							parseValid = false
							break
						}
					}
					if parseValid {
						valid++
					}
				}
			}
		}

		Expect(float32(valid) / float32(total)).To(BeNumerically(">=", 0.1))
	})
})

func BenchmarkParse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = address.Parse("331 Heather Hill Dr, Gibsonia, PA 15044")
	}
}
