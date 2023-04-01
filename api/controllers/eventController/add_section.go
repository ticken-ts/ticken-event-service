package eventController

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"ticken-event-service/api/mappers"
	"ticken-event-service/api/res"
	"ticken-event-service/security/jwt"
)

type createSectionPayload struct {
	Name         string  `json:"name"`
	TicketPrice  float64 `json:"ticket_price"`
	TotalTickets int     `json:"total_tickets"`
}

func (controller *EventController) AddSection(c *gin.Context) {
	var payload createSectionPayload

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

	if err = c.BindJSON(&payload); err != nil {
		c.Abort()
		return
	}

	section, err := controller.serviceProvider.GetEventManager().AddSection(
		organizerID,
		organizationID,
		eventID,
		payload.Name,
		payload.TotalTickets,
		payload.TicketPrice,
	)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, res.Success{
		Message: "section added successfully",
		Data:    mappers.MapSectionToDTO(section),
	})
}
