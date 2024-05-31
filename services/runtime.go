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
	timeout time.Duration
	vms     *runtime.Pool
}

var (
	ErrVMTimeout = errors.New("vm timed out")
)

func NewRuntime(
	client *sql.DB,
	timeout time.Duration,
) *Runtime {
	return &Runtime{
		vms:     runtime.NewPool(client, timeout),
		timeout: timeout,
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

	timer := time.AfterFunc(r.timeout, func() {
		jsRuntime.Interrupt(ErrVMTimeout)
	})
	defer timer.Stop()

	program, err := goja.Compile(
		"main.js",
		fmt.Sprintf(`{(function() {%s}).apply(undefined)}`, source),
		true,
	)
	if err != nil {
		return nil, fmt.Errorf("could not compile: %w", err)
	}

	returnValue, err := jsRuntime.RunProgram(program)
	if err != nil {
		defer jsRuntime.ClearInterrupt()

		return nil, fmt.Errorf("could not run program: %w", err)
	}

	return returnValue, nil
}
