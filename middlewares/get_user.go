package middlewares

import (
	"github.com/gin-gonic/gin"
	"ticken-event-service/services"
)

func GetUserMiddleware(services services.Provider) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
	}
}
