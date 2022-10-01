package middlewares

import (
	"github.com/gin-gonic/gin"
	"ticken-event-service/services"
)

func GetUserMiddleware(services services.Provider) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		token := ctx.GetHeader("token")

		userId, err := services.GetUserManager().GetUserIdFromToken(token)
		if err != nil {
			ctx.String(401, err.Error())
			ctx.Abort()
		}

		ctx.Set("userId", userId)

		ctx.Next()
	}
}
