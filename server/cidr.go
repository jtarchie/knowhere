package server

import (
	"fmt"
	"net"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/libp2p/go-cidranger"
)

func CIDRAllow(cidrs ...string) func(echo.HandlerFunc) echo.HandlerFunc {
	ranger := cidranger.NewPCTrieRanger()

	for _, cidr := range cidrs {
		_, network, err := net.ParseCIDR(cidr)
		if err != nil {
			panic(fmt.Sprintf("the CIDR (%s) is invalid: %s", cidr, err.Error()))
		}

		_ = ranger.Insert(cidranger.NewBasicRangerEntry(*network))

	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			contains, _ := ranger.Contains(net.ParseIP(c.RealIP()))
			if contains {
				return next(c)
			}

			return c.NoContent(http.StatusForbidden)
		}
	}
}
