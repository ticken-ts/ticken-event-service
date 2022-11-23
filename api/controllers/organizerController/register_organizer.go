package organizerController

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"ticken-event-service/api/mappers"
	"ticken-event-service/api/security"
	"ticken-event-service/utils"
)

func (controller *OrganizerController) RegisterOrganizer(c *gin.Context) {
	jwt := c.MustGet("jwt").(*security.JWT)

	organizerManager := controller.serviceProvider.GetOrgManager()

	organizer, err := organizerManager.RegisterOrganizer(jwt.Subject, jwt.Email, jwt.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
		c.Abort()
		return
	}

	organizerDTO := mappers.MapOrganizerToOrganizerDTO(organizer)

	c.JSON(http.StatusCreated, utils.HttpResponse{Data: organizerDTO})
}
