// nolint
package address

import "regexp"

// source: https://github.com/Senzing/libpostal-data/blob/main/files/tests/v1.1.0/test_data.csv
var addressParsers = []*regexp.Regexp{
	regexp.MustCompile(`^(?<house>.+)\s+(?<house_number>\d+)\s+(?<road>[^,]+),\s+(?<city_district>[^,]+),\s+(?<state>[^,]+),\s+(?<postcode>\d+)$`),
	regexp.MustCompile(`^(?<house>[^,]+),\s+(?<house_number>\d+)\s+(?<road>[^,]+),\s+(?<level>[^,]+),\s+(?<city_district>.+)\s+(?<postcode>\d+)$`),
	regexp.MustCompile(`^(?<house_number>\d+)\s+(?<road>.+)\s+(?<level>\d+),\s+(?<city>.+)\s+(?<state>.+)\s+(?<postcode>\d+),\s+(?<country>.+)$`),
	regexp.MustCompile(`^(?<house_number>\d+)\s+(?<road>.+)\s+(?<unit>\d+),\s+(?<city>.+)\s+(?<state>.+)\s+(?<postcode>\d+),\s+(?<country>.+)$`),
	regexp.MustCompile(`^(?<house>.+)\s+(?<house_number>\d+)\s+(?<road>.+)\s+(?<unit>\d+),\s+(?<city>.+)\s+(?<state>.+)\s+(?<postcode>\d+)$`),
	regexp.MustCompile(`^(?<house_number>\d+)\s+(?<road>[^,]+),\s+(?<city>.+)\s+(?<state>.+)\s+(?<postcode>[^,]+),\s+(?<country>.+)$`),
	regexp.MustCompile(`^(?<house_number>\d+)\s+(?<road>.+)\s+(?<level>\d+),\s+(?<city>.+)\s+(?<postcode>\d+),\s+(?<country>.+)$`),
	regexp.MustCompile(`^(?<house_number>\d+)\s+(?<road>[^,]+),\s+(?<suburb>.+)\s+(?<city>.+)\s+(?<state>.+)\s+(?<postcode>\d+)$`),
	regexp.MustCompile(`^(?<house_number>\d+)\s+(?<road>.+)\s+(?<unit>[^,]+),\s+(?<city>.+)\s+(?<state>.+)\s+(?<postcode>\d+)$`),
	regexp.MustCompile(`^(?<house_number>\d+)\s+(?<road>.+)\s+(?<level>\d+),\s+(?<city>.+)\s+(?<state>.+)\s+(?<postcode>\d+)$`),
	regexp.MustCompile(`^(?<house_number>\d+)\s+(?<road>[^,]+),\s+(?<city>.+)\s+(?<postcode>\d+),\s+(?<country>.+)$`),
	regexp.MustCompile(`^(?<house_number>\d+)\s+(?<road>[^,]+),\s+(?<city>.+)\s+(?<state>.+)\s+(?<postcode>\d+)$`),
	regexp.MustCompile(`^(?<house_number>\d+)\s+(?<road>.+)\s+(?<unit>[^,]+),\s+(?<city>.+)\s+(?<postcode>\d+)$`),
	regexp.MustCompile(`^(?<house_number>\d+)\s+(?<road>.+)\s+(?<level>\d+),\s+(?<city>.+)\s+(?<postcode>\d+)$`),
	regexp.MustCompile(`^(?<house_number>\d+)\s+(?<road>.+)\s+(?<unit>[^,]+),\s+(?<city>.+)\s+(?<state>.+)$`),
	regexp.MustCompile(`^(?<house_number>\d+)\s+(?<road>.+)\s+(?<level>\d+),\s+(?<city>.+)\s+(?<state>.+)$`),
	regexp.MustCompile(`^(?<house_number>\d+)\s+(?<road>[^,]+),\s+(?<city>.+)\s+(?<postcode>.+)$`),
	regexp.MustCompile(`^(?<house_number>\d+)\s+(?<road>[^,]+),\s+(?<city>.+)\s+(?<state>.+)$`),
	regexp.MustCompile(`^(?<road>[^,]+),\s+(?<suburb>[^,]+),\s+(?<city_district>.+)$`),
	regexp.MustCompile(`^(?<house_number>\d+)\s+(?<road>[^,]+),\s+(?<city>.+)$`),
	regexp.MustCompile(`^(?<road>[^,]+),\s+(?<city>.+)\s+(?<state>.+)$`),
	regexp.MustCompile(`^(?<city_district>[^,]+),\s+(?<state>.+)$`),
	regexp.MustCompile(`^(?<house_number>\d+)\s+(?<road>.+)$`),
	regexp.MustCompile(`^(?<city>[^,]+),\s+(?<state>.+)$`),
	regexp.MustCompile(`^(?<house>.+)\s+(?<city>.+)$`),
}
