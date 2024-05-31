package runtime_test

import (
	"testing"
	"time"

	"github.com/jtarchie/knowhere/runtime"
)

func BenchmarkPool(b *testing.B) {
	pool := runtime.NewPool(nil, time.Second)

	b.ResetTimer()

	b.RunParallel(func(p *testing.PB) {
		vm, err := pool.Get()
		if err != nil {
			b.Fatalf("could not get VM: %s", err)
		}
		defer pool.Put(vm)

		for p.Next() {
			_, err = vm.RunString(`{(function() {return "Hello, World"}).apply(undefined)}`)
			if err != nil {
				b.Fatalf("could not execute VM: %s", err)
			}
		}
	})
}
