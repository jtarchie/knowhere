package address_test

import (
	"log/slog"
	"os"
	"strings"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/jtarchie/knowhere/address"
	"github.com/recursionpharma/go-csv-map"
	"github.com/samber/lo/mutable"

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
		mutable.Shuffle(records)
		for _, record := range records {
			if record["country_code"] == "us" {
				total++
				fullAddress := record["full_address"]
				parsedAddress, ableToParse := address.Parse(fullAddress, false)

				if ableToParse {
					parseMatches := true

					for key, value := range parsedAddress {
						if strings.ToLower(value) != record[key] {
							slog.Debug("matches", "full", fullAddress, "key", key, "parsed", strings.ToLower(value), "record", record[key])
							parseMatches = false
							break
						}
					}

					if parseMatches {
						valid++
					}
				}
			}
		}

		Expect(float32(valid) / float32(total)).To(BeNumerically(">=", 0.75))
	})
})

func BenchmarkParse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parts, ok := address.Parse("331 Heather Hill Dr, Gibsonia, PA 15044", false)
		if !ok {
			b.Fatalf("could not because %#v", parts)
		}

		assert.Equal(b, parts["house_number"], "331")
		assert.Equal(b, parts["road"], "Heather Hill Dr")
		assert.Equal(b, parts["city"], "Gibsonia")
		assert.Equal(b, parts["state"], "PA")
		assert.Equal(b, parts["postcode"], "15044")
	}
}

func BenchmarkParseCleanup(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parts, ok := address.Parse("331 Heather Hill Dr, Gibsonia, PA 15044", true)
		if !ok {
			b.Fatalf("could not because %#v", parts)
		}

		assert.Equal(b, parts["house_number"], "331")
		assert.Equal(b, parts["road"], "Heather Hill Drive")
		assert.Equal(b, parts["city"], "Gibsonia")
		assert.Equal(b, parts["state"], "Pennsylvania")
		assert.Equal(b, parts["postcode"], "15044")
	}
}
