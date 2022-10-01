package api

import (
	"github.com/gin-gonic/gin"
	"ticken-event-service/infra"
)

type Controller interface {
	Setup(router infra.Router)
}

type Middleware = gin.HandlerFunc
