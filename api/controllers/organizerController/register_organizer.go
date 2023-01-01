package organizerController

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"ticken-event-service/api/mappers"
	"ticken-event-service/api/security"
	"ticken-event-service/utils"
)

type registerOrganizerPayload struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

func (controller *OrganizerController) CreateOrganizer(c *gin.Context) {
	jwt := c.MustGet("jwt").(*security.JWT)

	var payload registerOrganizerPayload

	err := c.BindJSON(&payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
		c.Abort()
		return
	}

	organizer, err := controller.serviceProvider.GetOrganizerManager().RegisterOrganizer(
		jwt.Subject,
		payload.Firstname,
		payload.Lastname,
		jwt.Username,
		jwt.Email,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
		c.Abort()
		return
	}

	eventsDTO := mappers.MapOrganizerToOrganizerDTO(organizer)

	c.JSON(http.StatusOK, utils.HttpResponse{Data: eventsDTO})
}
