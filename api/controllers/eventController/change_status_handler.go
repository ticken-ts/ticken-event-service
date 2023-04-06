package eventController

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"ticken-event-service/api/mappers"
	"ticken-event-service/api/res"
	"ticken-event-service/models"
	"ticken-event-service/security/jwt"
)

type statusChangePayload struct {
	NextStatus models.EventStatus `json:"next_status"`
}

func (controller *EventController) ChangeStatusHandler(c *gin.Context) {
	organizerID := c.MustGet("jwt").(*jwt.Token).Subject

	var payload statusChangePayload

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

	if err := c.BindJSON(&payload); err != nil {
		c.Abort()
		return
	}

	var event *models.Event
	var statusChangeErr error

	switch payload.NextStatus {
	case models.EventStatusOnSale:
		event, statusChangeErr = controller.ServiceProvider.GetEventManager().StartSale(
			eventID,
			organizerID,
			organizationID,
		)
	case models.EventStatusRunning:
		event, statusChangeErr = controller.ServiceProvider.GetEventManager().StartEvent(
			eventID,
			organizerID,
			organizationID,
		)
	case models.EventStatusFinished:
		event, statusChangeErr = controller.ServiceProvider.GetEventManager().FinishEvent(
			eventID,
			organizerID,
			organizationID,
		)
	default:
		statusChangeErr = fmt.Errorf("status change to %s is not supported", payload.NextStatus)
	}

	if statusChangeErr != nil {
		c.Error(statusChangeErr)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, res.Success{
		Message: "event set on sale successfully",
		Data:    mappers.MapEventToEventDTO(event),
	})
}
