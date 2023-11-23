package marshal

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/jtarchie/knowhere/query"
)

func Results(writer io.Writer, results []query.Result) error {
	if len(results) == 0 {
		_, _ = io.WriteString(writer, "[]")

		return nil
	}

	encoder := json.NewEncoder(writer)

	err := encoder.Encode(results)
	if err != nil {
		return fmt.Errorf("could not encode JSON: %w", err)
	}

	return nil
}
