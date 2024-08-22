package server

import (
	"bytes"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/georgysavva/scany/v2/sqlscan"
	"github.com/labstack/echo/v4"
)

type Area struct {
	Name     string  `db:"name"`
	FullName string  `db:"full_name"`
	MinLat   float64 `db:"minLat"`
	MaxLat   float64 `db:"maxLat"`
	MinLon   float64 `db:"minLon"`
	MaxLon   float64 `db:"maxLon"`
}

func (p *Area) MarshalJSON() ([]byte, error) {
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

func areas(client *sql.DB) func(echo.Context) error {
	return func(ctx echo.Context) error {
		var areas []Area

		err := sqlscan.Select(
			ctx.Request().Context(),
			client,
			&areas, `
			SELECT
				name, full_name, minLat, maxLat, minLon, maxLon
			FROM
				areas
		`)
		if err != nil {
			slog.Error("areas.error", slog.String("error", err.Error()))

			return response(ctx, http.StatusBadRequest, map[string]string{
				"error": "Results could not be processed",
			})
		}

		ctx.Response().Header().Set("Cache-Control", "public, max-age=1800")
		ctx.Response().Header().Set("Expires", time.Now().Add(30*time.Minute).Format(http.TimeFormat))

		return response(ctx, http.StatusOK, map[string]interface{}{
			"areas": areas,
		})
	}
}
