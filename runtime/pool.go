package runtime

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	_ "embed"

	"github.com/dop251/goja"
	"github.com/jtarchie/knowhere/query"
	"github.com/samber/lo"
)

type Pool struct {
	pool sync.Pool
}

var (
	ErrVMUnavailable = errors.New("could not get vm")
)

func NewPool(client *sql.DB, timeout time.Duration) *Pool {
	return &Pool{
		pool: sync.Pool{
			New: func() any {
				vm := goja.New() //nolint: varnamelen

				vm.SetFieldNameMapper(&tagFieldNameMapper{
					cache: map[string]string{},
				})

				err := vm.Set("rtree", func() *RTree {
					return &RTree{}
				})
				if err != nil {
					return fmt.Errorf("could not setup `rtree` VM: %w", err)
				}

				err = vm.Set("execute", func(qs string) []WrappedResult {
					ctx, cancel := context.WithTimeout(context.TODO(), timeout)
					defer cancel()

					results, err := query.Execute(ctx, client, qs, query.ToIndexedSQL)
					if err != nil {
						slog.Error("execute.failed", "query", qs, "err", err.Error())
						vm.Interrupt(fmt.Sprintf("could not execute query: %q", qs))
					}

					return lo.Map(results, func(result query.Result, _ int) WrappedResult {
						return WrappedResult{result}
					})
				})
				if err != nil {
					return fmt.Errorf("could not setup `execute` VM: %w", err)
				}

				err = vm.Set("assert", &Assertion{VM: vm})
				if err != nil {
					return fmt.Errorf("could not setup `assert` VM: %w", err)
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
