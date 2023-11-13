package marshal

import "fmt"

func Tags(tags map[string]string) string {
	if len(tags) == 0 {
		return "{}"
	}

	payload := "{"

	count := 0

	for key, value := range tags {
		payload += fmt.Sprintf("%q:%q", key, value)
		if count < len(tags)-1 {
			payload += ","
		}
		count++
	}

	payload += "}"

	return payload
}
