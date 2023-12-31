package errors

import (
	"fmt"
	"runtime"

	"github.com/labstack/echo/v4"

	"github.com/smiletrl/micro_ecommerce/pkg/constants"
	"github.com/smiletrl/micro_ecommerce/pkg/logger"
)

func Middleware(logger logger.Provider) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			defer func() {
				if r := recover(); r != nil {
					err, ok := r.(error)
					if !ok {
						err = fmt.Errorf("%v", r)
					}
					stack := make([]byte, 4<<10)
					length := runtime.Stack(stack, true)
					msg := fmt.Sprintf("[PANIC RECOVER] %v %s\n", err, stack[:length])
					logger.Errorw(msg)

					c.Error(err)
				}
			}()

			// set logger.
			c.Set(constants.Logger, logger)
			return next(c)
		}
	}
}
