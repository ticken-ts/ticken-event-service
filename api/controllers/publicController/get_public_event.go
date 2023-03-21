package publicController

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"ticken-event-service/api/mappers"
	"ticken-event-service/utils"
)

func (controller *PublicController) GetPublicEvent(c *gin.Context) {
	eventID, err := uuid.Parse(c.Param("eventID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
		c.Abort()
		return
	}

	manager := controller.serviceProvider.GetEventManager()
	event, err := manager.GetPublicEvent(eventID)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
		c.Abort()
		return
	}

	eventDTO := mappers.MapEventToEventDTO(event)
	c.JSON(http.StatusOK, utils.HttpResponse{Data: eventDTO})

}
