package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"ticken-event-service/api/res"
	"ticken-event-service/log"
	"ticken-event-service/tickenerr"
)

var statusCodesByError = map[uint32]int{}

type ErrorMiddleware struct {
}

func NewErrorMiddleware() *ErrorMiddleware {
	return new(ErrorMiddleware)
}

func (middleware *ErrorMiddleware) Setup(router gin.IRouter) {
	router.Use(middleware.ErrorHandler())
}

func (middleware *ErrorMiddleware) ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		for _, err := range c.Errors {
			switch v := err.Err.(type) {
			case tickenerr.TickenError:
				if v.UnderlyingErr != nil {
					log.TickenLogger.Error().Msg(v.UnderlyingErr.Error())
				}
				c.JSON(getStatusCode(v.Code), res.Error{Code: v.Code, Message: v.Message})

			default:
				c.JSON(http.StatusInternalServerError, res.Error{Code: 0, Message: "An error occurred"})
			}
		}
	}
}

func getStatusCode(errCode uint32) int {
	statusCode, ok := statusCodesByError[errCode]
	if !ok {
		return http.StatusInternalServerError
	}
	return statusCode
}
