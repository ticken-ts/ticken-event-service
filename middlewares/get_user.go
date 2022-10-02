package middlewares

import (
	"github.com/gin-gonic/gin"
	"ticken-event-service/api/errors"
	"ticken-event-service/services"
)

func GetUserMiddleware(services services.Provider) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		token := ctx.GetHeader("token")

		userId, err := services.GetUserManager().GetUserIdFromToken(token)
		if err != nil {
			apiError := errors.GetApiError(err)
			ctx.String(apiError.HttpCode, apiError.Message)
			ctx.Abort()
		}

		ctx.Set("userId", userId)

		ctx.Next()
	}
}
