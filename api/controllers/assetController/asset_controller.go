package assetController

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"ticken-event-service/services"
)

type AssetController struct {
	validator       *validator.Validate
	serviceProvider services.IProvider
}

func New(serviceProvider services.IProvider) *AssetController {
	controller := new(AssetController)
	controller.validator = validator.New()
	controller.serviceProvider = serviceProvider
	return controller
}

func (controller *AssetController) Setup(router gin.IRouter) {
	group := router.Group("/assets")
	group.GET("/:assetID", controller.GetAsset)
}
