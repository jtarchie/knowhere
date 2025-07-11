package services_test

import (
	"database/sql"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	"github.com/jtarchie/knowhere/services"
	_ "github.com/mattn/go-sqlite3"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Converter", func() {
	value := func(dbPath string, query string, result any) {
		client, err := sqlx.Open("sqlite3", dbPath)
		Expect(err).NotTo(HaveOccurred())

		err = client.Get(result, query)
		Expect(err).NotTo(HaveOccurred())
	}

	It("limits tags to a defined set", func() {
		buildDir, err := os.MkdirTemp("", "")
		Expect(err).NotTo(HaveOccurred())

		dbPath := filepath.Join(buildDir, "test.db")
		builder := services.NewConverter("../fixtures/sample.pbf", dbPath, "united_states", []string{"name"}, true, true)

		err = builder.Execute()
		Expect(err).NotTo(HaveOccurred())

		var tagCount int64
		value(dbPath, "SELECT COUNT(DISTINCT json_each.key) FROM united_states_entries, json_each(united_states_entries.tags);", &tagCount)
		Expect(tagCount).To(BeEquivalentTo(2)) // added geohash
	})

	It("puts nodes, ways, and relations into the entries", func() {
		buildDir, err := os.MkdirTemp("", "")
		Expect(err).NotTo(HaveOccurred())

		dbPath := filepath.Join(buildDir, "test.db")
		builder := services.NewConverter("../fixtures/sample.pbf", dbPath, "united_states", []string{"*"}, true, true)

		err = builder.Execute()
		Expect(err).NotTo(HaveOccurred())

		var result int64

		value(dbPath, "SELECT COUNT(*) FROM united_states_entries", &result)
		Expect(result).To(BeEquivalentTo(59))

		value(dbPath, "SELECT COUNT(*) FROM united_states_entries WHERE osm_type = 1", &result)
		Expect(result).To(BeEquivalentTo(10))

		value(dbPath, "SELECT COUNT(*) FROM united_states_entries WHERE osm_type = 2", &result)
		Expect(result).To(BeEquivalentTo(44))

		value(dbPath, "SELECT COUNT(*) FROM united_states_entries WHERE osm_type = 3", &result)
		Expect(result).To(BeEquivalentTo(5))

		value(dbPath, "SELECT COUNT(*) FROM united_states_entries WHERE osm_type = 2 AND tags <> '{}'", &result)
		Expect(result).To(BeEquivalentTo(44))

		value(dbPath, "SELECT COUNT(*) FROM united_states_entries WHERE osm_type = 2", &result)
		Expect(result).To(BeEquivalentTo(44))

		value(dbPath, "SELECT COUNT(*) FROM united_states_entries WHERE osm_type = 3 AND tags <> '{}'", &result)
		Expect(result).To(BeEquivalentTo(5))

		value(dbPath, "SELECT COUNT(*) FROM united_states_entries WHERE osm_type = 3", &result)
		Expect(result).To(BeEquivalentTo(5))

		// checking the id of full text search matches the id in the entries table
		var searchID, wayID int64
		value(dbPath, "SELECT MIN(rowid) FROM united_states_search WHERE tags MATCH 'Hatfield Tunnel' LIMIT 1", &searchID)
		value(dbPath, "SELECT id FROM united_states_entries WHERE tags->>'name' LIKE 'Hatfield Tunnel' AND osm_type = 2;", &wayID)
		Expect(searchID).To(BeEquivalentTo(wayID))

		var tagCount int64
		value(dbPath, "SELECT COUNT(DISTINCT json_each.key) FROM united_states_entries, json_each(united_states_entries.tags);", &tagCount)
		Expect(tagCount).To(BeEquivalentTo(47)) // HINT: geohash was added

		/*
			Napkin math for bounding box
			25365927: 51.7659279, -0.2326975
			691202858: 51.7663325, -0.2326806
			minLat, maxLat, minLon, maxLon
			51.7659279, 51.7663325,  -0.2326975,  -0.2326806
		*/
		var points struct {
			MinLat sql.NullFloat64 `db:"minLat"`
			MaxLat sql.NullFloat64 `db:"maxLat"`
			MinLon sql.NullFloat64 `db:"minLon"`
			MaxLon sql.NullFloat64 `db:"maxLon"`
		}
		value(dbPath, "SELECT minLat, maxLat, minLon, maxLon FROM united_states_entries WHERE id = 330;", &points)
		Expect(points.MinLat.Valid).To(BeTrue())
		Expect(points.MaxLat.Valid).To(BeTrue())
		Expect(points.MinLon.Valid).To(BeTrue())
		Expect(points.MaxLon.Valid).To(BeTrue())

		Expect(points.MinLat.Float64).To(BeNumerically("~", 51.76593))
		Expect(points.MaxLat.Float64).To(BeNumerically("~", 51.76633))
		Expect(points.MinLon.Float64).To(BeNumerically("~", -0.2327))
		Expect(points.MaxLon.Float64).To(BeNumerically("~", -0.23268))

		var area string
		value(dbPath, "SELECT name FROM areas", &area)
		Expect(area).To(Equal("united_states"))

		value(dbPath, "SELECT full_name FROM areas", &area)
		Expect(area).To(Equal("United States"))

		var fullCount int64
		value(dbPath, "SELECT COUNT(*) FROM united_states_entries WHERE minLat IS NULL", &fullCount)
		Expect(fullCount).To(BeEquivalentTo(0))
		value(dbPath, "SELECT COUNT(*) FROM united_states_entries WHERE minLon IS NULL", &fullCount)
		Expect(fullCount).To(BeEquivalentTo(0))
		value(dbPath, "SELECT COUNT(*) FROM united_states_entries WHERE maxLat IS NULL", &fullCount)
		Expect(fullCount).To(BeEquivalentTo(0))
		value(dbPath, "SELECT COUNT(*) FROM united_states_entries WHERE maxLon IS NULL", &fullCount)
		Expect(fullCount).To(BeEquivalentTo(0))
	})
})
