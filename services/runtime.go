package services

import (
	"database/sql"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/dop251/goja"
	"github.com/jtarchie/knowhere/query"
	"github.com/tidwall/geojson"
	"github.com/tidwall/geojson/geometry"
)

//go:embed turf.js
var turfJSSource string

type Runtime struct {
	vms sync.Pool
}

func NewRuntime(
	client *sql.DB,
) *Runtime {
	return &Runtime{
		vms: sync.Pool{
			New: func() any {
				vm := goja.New() //nolint: varnamelen

				vm.SetFieldNameMapper(goja.TagFieldNameMapper("js", true))

				_, err := vm.RunString(turfJSSource)
				if err != nil {
					return fmt.Errorf("could not warmup the VM: %w", err)
				}

				err = vm.Set("execute", func(qs string) any {
					results, err := query.Execute(client, qs, query.ToIndexedSQL)
					if err != nil {
						return fmt.Errorf("could not execute results: %w", err)
					}

					return results
				})
				if err != nil {
					return fmt.Errorf("could not setup `execute` VM: %w", err)
				}

				err = vm.Set("assertGeoJSON", func(payload any) bool {
					contents, err := json.Marshal(payload)
					if err != nil {
						return false
					}

					_, err = geojson.Parse(string(contents), &geojson.ParseOptions{
						IndexChildren:     64,
						IndexGeometry:     64,
						IndexGeometryKind: geometry.QuadTree,
						RequireValid:      true,
						AllowSimplePoints: false,
						DisableCircleType: false,
						AllowRects:        false,
					})

					return err == nil
				})
				if err != nil {
					return fmt.Errorf("could not setup `assertGeoJSON` VM: %w", err)
				}

				return vm
			},
		},
	}
}

//nolint: ireturn
func (r *Runtime) Execute(
	source string,
) (goja.Value, error) {
	switch vm := r.vms.Get().(type) {
	case *goja.Runtime:
		defer r.vms.Put(vm)

		timer := time.AfterFunc(time.Second, func() {
			vm.Interrupt("halt")
		})
		defer timer.Stop()

		value, err := vm.RunString(fmt.Sprintf(`
			(function() {
				%s
			})()
		`, source))
		if err != nil {
			return nil, fmt.Errorf("could not run program: %w", err)
		}

		return value, nil
	case error:
		return nil, fmt.Errorf("could not get vm: %w", vm)
	default:
		return nil, errors.New("could get vm")
	}
}
