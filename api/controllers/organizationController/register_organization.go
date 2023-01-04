package organizationController

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"ticken-event-service/api/security"
	"ticken-event-service/utils"
)

type registerOrganizationPayload struct {
	Name string `json:"name"`
}

func (controller *OrganizationController) RegisterOrganization(c *gin.Context) {
	var payload registerOrganizationPayload

	jwt := c.MustGet("jwt").(*security.JWT)

	err := c.BindJSON(&payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
		c.Abort()
		return
	}

	organizationManager := controller.serviceProvider.GetOrganizationManager()

	organization, err := organizationManager.RegisterOrganization(payload.Name, jwt.Subject)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, utils.HttpResponse{Data: organization})
}
