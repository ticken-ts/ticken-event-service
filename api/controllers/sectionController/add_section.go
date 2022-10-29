package sectionController

import (
	"github.com/coreos/go-oidc"
	"github.com/gin-gonic/gin"
	"net/http"
	"ticken-event-service/api/mappers"
	"ticken-event-service/utils"
)

type createSectionPayload struct {
	Name         string `json:"section"`
	TotalTickets int    `json:"total_tickets"`
}

func (controller *SectionController) AddSection(c *gin.Context) {
	var payload createSectionPayload

	eventID := c.Param("eventID")
	userID := c.MustGet("jwt").(*oidc.IDToken).Subject

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

	section, err := eventManager.AddSection(userID, eventID, payload.Name, payload.TotalTickets)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
		c.Abort()
		return
	}

	sectionDTO := mappers.MapSectionToDTO(section)

	c.JSON(http.StatusOK, utils.HttpResponse{Data: sectionDTO})
}
