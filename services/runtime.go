package services

import (
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"time"

	"github.com/dop251/goja"
	"github.com/evanw/esbuild/pkg/api"
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

	result := api.Transform(source, api.TransformOptions{
		Loader:    api.LoaderTS,
		Format:    api.FormatCommonJS,
		Target:    api.ES2015,
		Sourcemap: api.SourceMapNone,
		Platform:  api.PlatformNeutral,
	})

	if len(result.Errors) > 0 {
		return nil, &goja.CompilerSyntaxError{
			CompilerError: goja.CompilerError{
				Message: result.Errors[0].Text,
			},
		}
	}

	program, err := goja.Compile(
		"main.js",
		"{(function() { const module = {}; "+string(result.Code)+"; return module.exports.payload;}).apply(undefined)}",
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
