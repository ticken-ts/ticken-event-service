package eventController

import (
	"github.com/coreos/go-oidc"
	"github.com/gin-gonic/gin"
	"net/http"
	"ticken-event-service/api/mappers"
	"ticken-event-service/utils"
	"time"
)

type createEventPayload struct {
	Name string `json:"section"`
	Date string `json:"date"`
}

func (controller *EventController) CreateEvent(c *gin.Context) {
	var payload createEventPayload

	// TODO! use this to submit the event in behaviour of user
	_ = c.MustGet("jwt").(*oidc.IDToken).Subject

	err := c.BindJSON(&payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
		c.Abort()
		return
	}

	// the only thing that we are going to validate
	// is the that we can bind the request to the struct
	// and the date format that we need to operate.
	// No further validation are going to be added, so all
	// validations are going to be performed on chain

	parsedDate, err := time.Parse(time.RFC3339, payload.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
		c.Abort()
		return
	}

	eventManager := controller.serviceProvider.GetEventManager()

	event, err := eventManager.CreateEvent(payload.Name, parsedDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
		c.Abort()
		return
	}

	eventDTO := mappers.MapEventToCreatedTicketDTO(event)

	c.JSON(http.StatusOK, utils.HttpResponse{Data: eventDTO})
}