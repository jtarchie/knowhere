package runtime_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/jtarchie/knowhere/runtime"
)

var _ = Describe("Assert", func() {
	It("asserts for equality", func() {
		pool := runtime.NewPool(nil, time.Second)

		vm, err := pool.Get()
		Expect(err).NotTo(HaveOccurred())
		defer pool.Put(vm)

		_, err = vm.RunString(`assert.eq(true, "true message")`)
		Expect(err).NotTo(HaveOccurred())

		_, err = vm.RunString(`assert.eq(false, "false message")`)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("false message"))
	})

	It("asserts valid GeoJSON", func() {
		pool := runtime.NewPool(nil, time.Second)

		vm, err := pool.Get()
		Expect(err).NotTo(HaveOccurred())
		defer pool.Put(vm)

		_, err = vm.RunString(`assert.geoJSON({})`)
		Expect(err).To(HaveOccurred())

		_, err = vm.RunString(`assert.geoJSON({
			"type": "FeatureCollection",
			"features": []
		})`)
		Expect(err).NotTo(HaveOccurred())
	})
})
