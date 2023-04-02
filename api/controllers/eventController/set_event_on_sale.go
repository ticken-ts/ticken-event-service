package eventController

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"ticken-event-service/api/mappers"
	"ticken-event-service/api/res"
	"ticken-event-service/security/jwt"
)

func (controller *EventController) SetEventOnSale(c *gin.Context) {
	organizerID := c.MustGet("jwt").(*jwt.Token).Subject

	organizationID, err := uuid.Parse(c.Param("organizationID"))
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	eventID, err := uuid.Parse(c.Param("eventID"))
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	event, err := controller.ServiceProvider.GetEventManager().SetEventOnSale(
		eventID,
		organizationID,
		organizerID,
	)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, res.Success{
		Message: "event set on sale successfully",
		Data:    mappers.MapEventToEventDTO(event),
	})
}
