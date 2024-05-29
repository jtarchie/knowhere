package services

import (
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"time"

	"github.com/dop251/goja"
	"github.com/jtarchie/knowhere/runtime"
)

type Runtime struct {
	vms *runtime.Pool
}

var (
	ErrVMTimeout = errors.New("vm timed out")
)

func NewRuntime(
	client *sql.DB,
) *Runtime {
	return &Runtime{
		vms: runtime.NewPool(client),
	}
}

// nolint: ireturn
func (r *Runtime) Execute(
	source string,
) (goja.Value, error) {
	jsRuntime, err := r.vms.Get()
	if err != nil {
		return nil, fmt.Errorf("could not get vm: %w", err)
	}
	defer r.vms.Put(jsRuntime)

	timer := time.AfterFunc(time.Second, func() {
		jsRuntime.Interrupt(ErrVMTimeout)
	})
	defer timer.Stop()

	returnValue, err := jsRuntime.RunString(fmt.Sprintf(`
			(function() {
				%s
			})()
		`, source))
	if err != nil {
		return nil, fmt.Errorf("could not run program: %w", err)
	}

	return returnValue, nil
}
