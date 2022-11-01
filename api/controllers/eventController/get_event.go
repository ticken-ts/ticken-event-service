package eventController

import (
	"github.com/coreos/go-oidc"
	"github.com/gin-gonic/gin"
	"net/http"
	"ticken-event-service/api/mappers"
	"ticken-event-service/utils"
)

func (controller *EventController) GetEvent(c *gin.Context) {
	userID := c.MustGet("jwt").(*oidc.IDToken).Subject
	eventId := c.Param("eventId")

	event, err := controller.serviceProvider.GetEventManager().GetEvent(eventId, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
		c.Abort()
		return
	}

	eventDTO := mappers.MapEventToEventDTO(event)

	c.JSON(http.StatusOK, utils.HttpResponse{Data: eventDTO})
}
