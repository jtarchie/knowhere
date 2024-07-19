// nolint
package address

import "regexp"

// source: https://github.com/Senzing/libpostal-data/blob/main/files/tests/v1.1.0/test_data.csv
var addressParsers = []*regexp.Regexp{
	regexp.MustCompile(`^(?<house_number>\d+)\s+(?<road>[^,]+),\s+(?<city>.+)\s+(?<state>.+)\s+(?<postcode>\d+)$`),
	regexp.MustCompile(`^(?<house_number>\d+)\s+(?<road>[^,]+),\s+(?<city>.+)\s+(?<postcode>.+)$`),
	regexp.MustCompile(`^(?<house_number>\d+)\s+(?<road>[^,]+),\s+(?<city>.+)\s+(?<state>.+)\s+(?<postcode>[^,]+),\s+(?<country>.+)$`),
	regexp.MustCompile(`^(?<house_number>\d+)\s+(?<road>.+)\s+(?<unit>[^,]+),\s+(?<city>.+)\s+(?<state>.+)\s+(?<postcode>\d+)$`),
	regexp.MustCompile(`^(?<house_number>\d+)\s+(?<road>[^,]+),\s+(?<city>.+)\s+(?<postcode>\d+),\s+(?<country>.+)$`),
	regexp.MustCompile(`^(?<house_number>\d+)\s+(?<road>[^,]+),\s+(?<city>.+)\s+(?<state>.+)$`),
	regexp.MustCompile(`^(?<house_number>\d+)\s+(?<road>.+)\s+(?<unit>[^,]+),\s+(?<city>.+)\s+(?<postcode>\d+)$`),
	regexp.MustCompile(`^(?<house_number>\d+)\s+(?<road>.+)\s+(?<level>\d+),\s+(?<city>.+)\s+(?<state>.+)\s+(?<postcode>\d+)$`),
	regexp.MustCompile(`^(?<house_number>\d+)\s+(?<road>[^,]+),\s+(?<city>.+)$`),
	regexp.MustCompile(`^(?<house_number>\d+)\s+(?<road>.+)\s+(?<level>\d+),\s+(?<city>.+)\s+(?<postcode>\d+)$`),
}
