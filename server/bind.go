package server

import (
	"errors"
	"fmt"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/tinylib/msgp/msgp"
)

const msgpackContentType = "application/msgpack"

func response(context echo.Context, status int, payload interface{}) error {
	accept := context.Request().Header.Get(echo.HeaderAccept)

	if strings.Contains(accept, msgpackContentType) {
		encodable, ok := payload.(msgp.Encodable)
		if !ok {
			return errors.New("could not convert payload to msgpack")
		}

		context.Response().Header().Set(echo.HeaderContentType, msgpackContentType)
		context.Response().WriteHeader(status)

		err := msgp.Encode(context.Response(), encodable)
		if err != nil {
			return fmt.Errorf("could not marshal payload to msgpack: %w", err)
		}

		return nil
	}

	return context.JSON(status, payload)
}
