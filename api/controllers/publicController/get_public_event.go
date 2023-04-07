package publicController

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"ticken-event-service/api/mappers"
	"ticken-event-service/api/res"
)

func (controller *PublicController) GetPublicEvent(c *gin.Context) {
	eventID, err := uuid.Parse(c.Param("eventID"))
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	event, err := controller.serviceProvider.GetEventManager().GetPublicEvent(eventID)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, res.Success{
		Message: "event found successfully",
		Data:    mappers.MapEventToEventDTO(event),
	})
}
