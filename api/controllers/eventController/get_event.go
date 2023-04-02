package eventController

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"ticken-event-service/api/mappers"
	"ticken-event-service/api/res"
	"ticken-event-service/security/jwt"
	"ticken-event-service/utils"
)

func (controller *EventController) GetEvent(c *gin.Context) {
	organizerID := c.MustGet("jwt").(*jwt.Token).Subject

	eventID, err := uuid.Parse(c.Param("eventID"))
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	organizationID, err := uuid.Parse(c.Param("organizationID"))
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	event, err := controller.ServiceProvider.GetEventManager().GetEvent(
		eventID,
		organizerID,
		organizationID,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, res.Success{
		Message: "event retrieved successfully",
		Data:    mappers.MapEventToEventDTO(event),
	})
}
