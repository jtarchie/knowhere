package runtime

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/dop251/goja"
	"github.com/goccy/go-json"
	"github.com/tidwall/geojson"
	"github.com/tidwall/geojson/geometry"
)

type Assertion struct {
	vm   *goja.Runtime
	time time.Time
}

func (a *Assertion) Stab(msg string) {
	if a.time.IsZero() {
		a.time = time.Now()
	}

	slog.Info("stab", slog.String("msg", msg), slog.Duration("time", time.Since(a.time)))
	a.time = time.Now()
}

func (a *Assertion) Eq(value bool, msg string) {
	if !value {
		a.vm.Interrupt(fmt.Sprintf("assertion failed: %s", msg))
	}
}

func (a *Assertion) GeoJSON(payload any) {
	contents, err := json.Marshal(payload)
	if err != nil {
		a.vm.Interrupt("geojson payload is not JSON")

		return
	}

	_, err = geojson.Parse(string(contents), &geojson.ParseOptions{
		IndexGeometryKind: geometry.None,
		RequireValid:      true,
	})
	if err != nil {
		a.vm.Interrupt(fmt.Sprintf("assert of geojson failed: %s", err))
	}
}
