package server

import (
	"fmt"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/vmihailenco/msgpack/v5"
)

const msgpackContentType = "application/msgpack"

func response(context echo.Context, status int, payload interface{}) error {
	accept := context.Request().Header.Get(echo.HeaderAccept)

	if strings.Contains(accept, msgpackContentType) {
		context.Response().Header().Set(echo.HeaderContentType, msgpackContentType)
		context.Response().WriteHeader(status)

		encoder := msgpack.NewEncoder(context.Response())
		err := encoder.Encode(payload)
		if err != nil {
			return fmt.Errorf("could not marshal payload to msgpack: %w", err)
		}

		return nil
	}

	return context.JSON(status, payload)
}
