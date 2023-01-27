package eventController

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"ticken-event-service/api/mappers"
	"ticken-event-service/security/jwt"
	"ticken-event-service/utils"
)

func (controller *EventController) GetOrganizationEvents(c *gin.Context) {
	userID := c.MustGet("jwt").(*jwt.Token).Subject

	organizationID, err := uuid.Parse(c.Param("organizationID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
		c.Abort()
		return
	}

	eventManager := controller.serviceProvider.GetEventManager()

	events, err := eventManager.GetOrganizationEvents(userID, organizationID)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
		c.Abort()
		return
	}

	eventsDTO := mappers.MapEventListToDTO(events)

	c.JSON(http.StatusOK, utils.HttpResponse{Data: eventsDTO})
}
