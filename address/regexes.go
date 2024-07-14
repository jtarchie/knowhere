// nolint
package address

import "regexp"

// source: https://github.com/Senzing/libpostal-data/blob/main/files/tests/v1.1.0/test_data.csv
var addressParsers = []*regexp.Regexp{
	regexp.MustCompile("^(?<house_number>\\d+)\\s+(?<road>.*),\\s+(?<city>.*)\\s+(?<state>\\w+)\\s+(?<postcode>.*),\\s+(?<country>\\w+)$"),
	regexp.MustCompile("^(?<house_number>\\d+)\\s+(?<road>.*)\\s+(?<unit>.*),\\s+(?<city>.*)\\s+(?<state>\\w+)\\s+(?<postcode>\\d+)$"),
	regexp.MustCompile("^(?<house_number>\\d+)\\s+(?<road>.*),\\s+(?<city>.*)\\s+(?<postcode>\\d+),\\s+(?<country>\\w+)$"),
	regexp.MustCompile("^(?<house_number>\\w+)\\s+(?<road>.*),\\s+(?<city>.*)\\s+(?<state>\\w+)\\s+(?<postcode>\\d+)$"),
	regexp.MustCompile("^(?<house_number>\\d+)\\s+(?<road>.*)\\s+(?<unit>\\w+),\\s+(?<city>.*)\\s+(?<postcode>\\d+)$"),
	regexp.MustCompile("^(?<house_number>\\d+)\\s+(?<road>.*),\\s+(?<city>.*)\\s+(?<postcode>\\d+)$"),
	regexp.MustCompile("^(?<house_number>\\d+)\\s+(?<road>.*),\\s+(?<city>.*)\\s+(?<state>\\w+)$"),
}
