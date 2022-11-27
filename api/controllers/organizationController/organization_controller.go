package organizationController

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"ticken-event-service/api/security"
	"ticken-event-service/services"
	"ticken-event-service/utils"
)

type OrganizationController struct {
	validator       *validator.Validate
	serviceProvider services.IProvider
}

func New(serviceProvider services.IProvider) *OrganizationController {
	controller := new(OrganizationController)
	controller.validator = validator.New()
	controller.serviceProvider = serviceProvider
	return controller
}

func (controller *OrganizationController) Setup(router gin.IRouter) {
	router.GET("/organizations/:organizationID/crypto", controller.GetOrganizationCrypto)
	router.POST("/organizations", controller.RegisterOrganization)
}

func (controller *OrganizationController) GetOrganizationCrypto(c *gin.Context) {
	organizationID := c.Param("organizationID")
	jwt := c.MustGet("jwt").(*security.JWT)

	organizationManager := controller.serviceProvider.GetOrganizationManager()

	zip, err := organizationManager.GetOrganizationCryptoZipped(jwt.Subject, organizationID)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
		c.Abort()
		return
	}

	c.Writer.Header().Set("Content-Type", "application/zip")
	c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.zip\"", "crypto"))
	_, _ = c.Writer.Write(zip)
}
