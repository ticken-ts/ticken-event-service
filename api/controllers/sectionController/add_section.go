package sectionController

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"ticken-event-service/api/mappers"
	"ticken-event-service/security/jwt"
	"ticken-event-service/utils"
)

type createSectionPayload struct {
	Name         string  `json:"name"`
	TicketPrice  float64 `json:"ticket_price"`
	TotalTickets int     `json:"total_tickets"`
}

func (controller *SectionController) AddSection(c *gin.Context) {
	var payload createSectionPayload

	userID := c.MustGet("jwt").(*jwt.Token).Subject

	eventID, err := uuid.Parse(c.Param("eventID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
		c.Abort()
		return
	}

	organizationID, err := uuid.Parse(c.Param("organizationID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
		c.Abort()
		return
	}

	if err = c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
		c.Abort()
		return
	}

	// the only thing that we are going to validate
	// is the that we can bind the request to the struct
	// No further validation are going to be added, so all
	// validations are going to be performed on chain

	eventManager := controller.serviceProvider.GetEventManager()

	section, err := eventManager.AddSection(
		userID,
		organizationID,
		eventID,
		payload.Name,
		payload.TotalTickets,
		payload.TicketPrice,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
		c.Abort()
		return
	}

	sectionDTO := mappers.MapSectionToDTO(section)

	c.JSON(http.StatusOK, utils.HttpResponse{Data: sectionDTO})
}
