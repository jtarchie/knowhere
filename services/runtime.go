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
	"github.com/samber/lo"
	"github.com/tidwall/geojson"
	"github.com/tidwall/geojson/geometry"
)

//go:embed turf.js
var turfJSSource string

type Runtime struct {
	vms sync.Pool
}

var (
	ErrVMUnavailable = errors.New("could not get vm")
	ErrVMTimeout     = errors.New("vm timed out")
)

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
						vm.Interrupt(fmt.Sprintf("could not execute query: %q", qs))
					}

					return lo.Map(results, func(result query.Result, _ int) WrappedResult {
						return WrappedResult{result}
					})
				})
				if err != nil {
					return fmt.Errorf("could not setup `execute` VM: %w", err)
				}

				err = vm.Set("assert", func(value bool) {
					if !value {
						vm.Interrupt("assertion failed")
					}
				})
				if err != nil {
					return fmt.Errorf("could not setup `assert` VM: %w", err)
				}

				err = vm.Set("assertGeoJSON", func(payload any) {
					contents, err := json.Marshal(payload)
					if err != nil {
						vm.Interrupt("geojson payload is not JSON")

						return
					}

					_, err = geojson.Parse(string(contents), &geojson.ParseOptions{
						IndexGeometryKind: geometry.QuadTree,
						RequireValid:      true,
					})
					if err != nil {
						vm.Interrupt("assert of geojson failed")

						return
					}
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
	switch value := r.vms.Get().(type) {
	case *goja.Runtime:
		defer r.vms.Put(value)

		timer := time.AfterFunc(time.Second, func() {
			value.Interrupt(ErrVMTimeout)
		})
		defer timer.Stop()

		returnValue, err := value.RunString(fmt.Sprintf(`
			(function() {
				%s
			})()
		`, source))
		if err != nil {
			return nil, fmt.Errorf("could not run program: %w", err)
		}

		return returnValue, nil
	case error:
		return nil, fmt.Errorf("could not get vm: %w", value)
	default:
		return nil, ErrVMUnavailable
	}
}
