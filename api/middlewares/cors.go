package middlewares

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type CorsMiddleware struct {
}

func NewCorsMiddleware() *CorsMiddleware {
	middleware := new(CorsMiddleware)
	return middleware
}

func (middleware *CorsMiddleware) Setup(router gin.IRouter) {
	router.Use(middleware.addCorsHeaders())
}

func (middleware *CorsMiddleware) addCorsHeaders() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	})
}
