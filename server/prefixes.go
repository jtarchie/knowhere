package server

import (
	"bytes"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/georgysavva/scany/v2/sqlscan"
	"github.com/labstack/echo/v4"
)

type Prefix struct {
	Name     string  `db:"name"`
	FullName string  `db:"full_name"`
	MinLat   float32 `db:"minLat"`
	MaxLat   float32 `db:"maxLat"`
	MinLon   float32 `db:"minLon"`
	MaxLon   float32 `db:"maxLon"`
}

func (p *Prefix) MarshalJSON() ([]byte, error) {
	buffer := &bytes.Buffer{}

	buffer.WriteString(`{"name":`)
	buffer.WriteString(fmt.Sprintf("%q", p.FullName))
	buffer.WriteString(`,"slug":`)
	buffer.WriteString(fmt.Sprintf("%q", p.Name))
	buffer.WriteString(`,"bounds":[[`)
	buffer.WriteString(strconv.FormatFloat(float64(p.MinLon), 'f', -1, 32)) // lng,lat
	buffer.WriteByte(',')
	buffer.WriteString(strconv.FormatFloat(float64(p.MinLat), 'f', -1, 32)) // lng,lat
	buffer.WriteString("],[")
	buffer.WriteString(strconv.FormatFloat(float64(p.MaxLon), 'f', -1, 32)) // lng,lat
	buffer.WriteByte(',')
	buffer.WriteString(strconv.FormatFloat(float64(p.MaxLat), 'f', -1, 32)) // lng,lat
	buffer.WriteString("]]}")

	return buffer.Bytes(), nil
}

func prefixes(client *sql.DB) func(echo.Context) error {
	return func(ctx echo.Context) error {
		var prefixes []Prefix

		err := sqlscan.Select(
			ctx.Request().Context(),
			client,
			&prefixes, `
			SELECT
				name, full_name, minLat, maxLat, minLon, maxLon
			FROM
				prefixes
		`)
		if err != nil {
			slog.Error("prefixes.error", slog.String("error", err.Error()))

			return ctx.JSON(http.StatusBadRequest, map[string]string{
				"error": "Results could not be processed",
			})
		}

		return ctx.JSON(http.StatusOK, map[string]interface{}{
			"prefixes": prefixes,
		})
	}
}
