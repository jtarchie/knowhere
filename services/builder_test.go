package services_test

import (
	"database/sql"
	"os"
	"path/filepath"

	"github.com/jtarchie/knowhere/services"
	_ "github.com/mattn/go-sqlite3"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Builder", func() {
	sql := func(dbPath string, query string) interface{} {
		client, err := sql.Open("sqlite3", dbPath)
		Expect(err).NotTo(HaveOccurred())

		row := client.QueryRow(query)
		Expect(row.Err()).NotTo(HaveOccurred())

		var result interface{}
		err = row.Scan(&result)
		Expect(err).NotTo(HaveOccurred())

		return result
	}

	It("puts nodes, ways, and relations into the entries", func() {
		buildDir, err := os.MkdirTemp("", "")
		Expect(err).NotTo(HaveOccurred())

		dbPath := filepath.Join(buildDir, "test.db")
		builder := services.NewBuilder("../fixtures/sample.pbf", dbPath)

		err = builder.Execute()
		Expect(err).NotTo(HaveOccurred())

		result := sql(dbPath, "SELECT COUNT(*) FROM entries")
		Expect(result).To(BeEquivalentTo(339))

		result = sql(dbPath, "SELECT COUNT(*) FROM entries WHERE osm_type = 'node'")
		Expect(result).To(BeEquivalentTo(290))

		result = sql(dbPath, "SELECT COUNT(*) FROM entries WHERE osm_type = 'way'")
		Expect(result).To(BeEquivalentTo(44))

		result = sql(dbPath, "SELECT COUNT(*) FROM entries WHERE osm_type = 'relation'")
		Expect(result).To(BeEquivalentTo(5))
	})
})
