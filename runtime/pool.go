package runtime

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	_ "embed"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
	"github.com/dop251/goja_nodejs/url"
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

				new(require.Registry).Enable(vm)
				url.Enable(vm)

				err := vm.Set("query", &Query{
					vm:      vm,
					timeout: timeout,
					client:  client,
				})
				if err != nil {
					return fmt.Errorf("could not setup `execute` VM: %w", err)
				}

				err = vm.Set("assert", &Assertion{vm: vm})
				if err != nil {
					return fmt.Errorf("could not setup `assert` VM: %w", err)
				}

				_ = vm.Set("colors", &Colors{})
				_ = vm.Set("geo", &Geo{})

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
