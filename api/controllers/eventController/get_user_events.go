package eventController

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"ticken-event-service/api/mappers"
	"ticken-event-service/api/security"
	"ticken-event-service/utils"
)

func (controller *EventController) GetOrganizationEvents(c *gin.Context) {
	userID := c.MustGet("jwt").(*security.JWT).Subject
	organizationID := c.Param("organizationID")

	events, err := controller.serviceProvider.GetEventManager().GetOrganizationEvents(userID, organizationID)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
		c.Abort()
		return
	}

	eventsDTO := mappers.MapEventListToEventListDTO(events)

	c.JSON(http.StatusOK, utils.HttpResponse{Data: eventsDTO})
}
