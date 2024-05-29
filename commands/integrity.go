package commands

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/jtarchie/knowhere/query"
	"github.com/jtarchie/sqlitezstd"
	"github.com/olekukonko/tablewriter"
	"github.com/samber/lo"
)

type Integrity struct {
	DB    string `help:"db filename to import data to" required:""`
	Query string `default:"nwr[name=Starbucks]"        help:"query to test against, defaults to all prefixes" required:""`
}

// nolint: all
func (i *Integrity) Run() error {
	connectionString := fmt.Sprintf("file:%s?_query_only=true&immutable=true&mode=ro", i.DB)

	if strings.Contains(i.DB, ".zst") {
		err := sqlitezstd.Init()
		if err != nil {
			return fmt.Errorf("could not load sqlite zstd vfs: %w", err)
		}

		connectionString += "&vfs=zstd"
	}

	client, err := sql.Open("sqlite3", connectionString)
	if err != nil {
		return fmt.Errorf("could not open database file: %w", err)
	}

	exactSQL, _ := query.ToExactSQL(i.Query)
	indexSQL, _ := query.ToIndexedSQL(i.Query)

	fmt.Println("exact query: ")
	fmt.Println(exactSQL)

	fmt.Println("index query: ")
	fmt.Println(indexSQL)

	slog.Info("querying")

	exactResults, err := query.Execute(context.TODO(), client, i.Query, query.ToExactSQL)
	if err != nil {
		return fmt.Errorf("could not get exact results: %w", err)
	}

	slog.Info("exact.done")

	indexResults, err := query.Execute(context.TODO(), client, i.Query, query.ToIndexedSQL)
	if err != nil {
		return fmt.Errorf("could not get index results: %w", err)
	}

	slog.Info("index.done")

	exactIDs := lo.Associate(exactResults, func(result query.Result) (int64, struct{}) {
		return result.ID, struct{}{}
	})

	indexIDs := lo.Associate(indexResults, func(result query.Result) (int64, struct{}) {
		return result.ID, struct{}{}
	})

	var tp, fp, fn int

	for id := range indexIDs {
		if _, ok := exactIDs[id]; ok {
			tp++
		} else {
			fp++
		}
	}

	for id := range exactIDs {
		if _, ok := indexIDs[id]; !ok {
			fn++
		}
	}

	precision := float64(tp) / float64(tp+fp)
	recall := float64(tp) / float64(tp+fn)
	f1Score := 2 * (precision * recall) / (precision + recall)

	data := [][]string{
		[]string{"True Position", fmt.Sprintf("%d", tp), "index that appears in exact"},
		[]string{"False Positives", fmt.Sprintf("%d", fp), "index that does not appear in exact"},
		[]string{"False Negatives", fmt.Sprintf("%d", fn), "exact that does not appear in index"},
		[]string{"Precision", fmt.Sprintf("%f", precision), ""},
		[]string{"Recall", fmt.Sprintf("%f", recall), ""},
		[]string{"F1", fmt.Sprintf("%f", f1Score), "Precision to recall (1 = exact results)"},
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Value", "Description"})

	for _, v := range data {
		table.Append(v)
	}
	table.Render() // Send output

	return nil
}
