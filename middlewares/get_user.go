package middlewares

import (
	"github.com/gin-gonic/gin"
	"ticken-event-service/services"
)

func GetUserMiddleware(services services.Provider) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		token := ctx.GetHeader("token")

		userId := services.GetUserManager().GetUserIdFromToken(token)

		ctx.Set("userId", userId)

		ctx.Next()
	}
}
