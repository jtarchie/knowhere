package server_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/labstack/echo/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/jtarchie/knowhere/server"
)

var _ = Describe("CIDR", func() {
	It("allows all IPs", func() {
		handler := echo.New()
		handler.Use(server.CIDRAllow("0.0.0.0/0"))

		server := httptest.NewServer(handler)
		defer server.Close()

		response, err := http.Get(server.URL)
		Expect(err).NotTo(HaveOccurred())

		Expect(response.StatusCode).To(Equal(http.StatusNotFound))
	})

	DescribeTable("with configuration", func(allowlist string, remoteIP string, status int) {
		handler := echo.New()

		request := httptest.NewRequest(http.MethodGet, "/", nil)
		request.RemoteAddr = remoteIP
		response := httptest.NewRecorder()
		context := handler.NewContext(request, response)
		middleware := server.CIDRAllow(allowlist)(func(c echo.Context) error {
			return c.NoContent(http.StatusOK)
		})

		err := middleware(context)
		Expect(err).NotTo(HaveOccurred())

		Expect(response.Result().StatusCode).To(Equal(status))
	},
		Entry("blocked", "123.123.123.123/32", "223.123.123.123:1234", http.StatusForbidden),
		Entry("not blocked", "223.123.123.123/32", "223.123.123.123:1234", http.StatusOK),
		Entry("not blocked cidr", "223.123.123.0/24", "223.123.123.123:1234", http.StatusOK),
	)
})
