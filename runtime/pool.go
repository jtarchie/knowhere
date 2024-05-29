package runtime

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	_ "embed"

	"github.com/dop251/goja"
	"github.com/jtarchie/knowhere/query"
	"github.com/samber/lo"
	"github.com/tidwall/geojson"
	"github.com/tidwall/geojson/geometry"
)

//go:embed turf.js
var turfJSSource string

type Pool struct {
	pool sync.Pool
}

var (
	ErrVMUnavailable = errors.New("could not get vm")
)

func NewPool(client *sql.DB) *Pool {
	return &Pool{
		pool: sync.Pool{
			New: func() any {
				vm := goja.New() //nolint: varnamelen

				vm.SetFieldNameMapper(goja.TagFieldNameMapper("js", true))

				_, err := vm.RunString(turfJSSource)
				if err != nil {
					return fmt.Errorf("could not warmup the VM: %w", err)
				}

				err = vm.Set("execute", func(qs string) any {
					ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
					defer cancel()

					results, err := query.Execute(ctx, client, qs, query.ToIndexedSQL)
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
						IndexGeometryKind: geometry.None,
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

func (p *Pool) Get() (*goja.Runtime, error) {
	switch value := p.pool.Get().(type) {
	case *goja.Runtime:
		return value, nil
	case error:
		return nil, fmt.Errorf("vm pool unavailable: %w", value)
	default:
		return nil, ErrVMUnavailable
	}
}

func (p *Pool) Put(runtime *goja.Runtime) {
	p.pool.Put(runtime)
}
