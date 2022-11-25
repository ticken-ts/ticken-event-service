package sectionController

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"ticken-event-service/api/mappers"
	"ticken-event-service/api/security"
	"ticken-event-service/utils"
)

type createSectionPayload struct {
	Name         string `json:"name"`
	TotalTickets int    `json:"total_tickets"`
}

func (controller *SectionController) AddSection(c *gin.Context) {
	var payload createSectionPayload

	eventID := c.Param("eventID")
	organizationID := c.Param("organizationID")
	userID := c.MustGet("jwt").(*security.JWT).Subject

	err := c.BindJSON(&payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
		c.Abort()
		return
	}

	// the only thing that we are going to validate
	// is the that we can bind the request to the struct
	// No further validation are going to be added, so all
	// validations are going to be performed on chain

	eventManager := controller.serviceProvider.GetEventManager()

	section, err := eventManager.AddSection(userID, organizationID, eventID, payload.Name, payload.TotalTickets)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
		c.Abort()
		return
	}

	sectionDTO := mappers.MapSectionToDTO(section)

	c.JSON(http.StatusOK, utils.HttpResponse{Data: sectionDTO})
}
