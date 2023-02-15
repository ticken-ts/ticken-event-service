package eventController

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"mime/multipart"
	"net/http"
	"ticken-event-service/api/mappers"
	"ticken-event-service/models"
	"ticken-event-service/security/jwt"
	"ticken-event-service/utils"
	"time"
)

type createEventPayload struct {
	Name        string       `json:"name" binding:"required"`
	Date        time.Time    `json:"date" binding:"required"`
	Description string       `json:"description"`
	PosterFile  *models.File `json:"poster"`
}

func (controller *EventController) CreateEvent(c *gin.Context) {
	userID := c.MustGet("jwt").(*jwt.Token).Subject

	organizationID, err := uuid.Parse(c.Param("organizationID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
		c.Abort()
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
		c.Abort()
		return
	}

	// validate the form
	payload, err := validateMultipartForm(form)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
		c.Abort()
		return
	}

	// the only thing that we are going to validate
	// is the that we can bind the request to the struct
	// No further validation are going to be added, so all
	// validations are going to be performed on chain

	if form.File["poster"] != nil {
		file, err := form.File["poster"][0].Open()
		if err != nil {
			c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
			c.Abort()
			return
		}
		bytes := make([]byte, form.File["poster"][0].Size)
		_, err = file.Read(bytes)
		if err != nil {
			c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
			c.Abort()
			return
		}
		payload.PosterFile = models.NewFile(&bytes, form.File["poster"][0].Header.Get("Content-Type"))
	}

	eventManager := controller.serviceProvider.GetEventManager()

	event, err := eventManager.CreateEvent(userID, organizationID, payload.Name, payload.Date, payload.Description, payload.PosterFile)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
		c.Abort()
		return
	}

	eventDTO := mappers.MapEventToEventDTO(event)

	c.JSON(http.StatusCreated, utils.HttpResponse{Data: eventDTO})
}

// validateMultipartForm will validate the form and return an error if the form is not valid
func validateMultipartForm(form *multipart.Form) (*createEventPayload, error) {
	if form.Value["name"] == nil {
		return nil, fmt.Errorf("name is required")
	}

	if form.Value["date"] == nil {
		return nil, fmt.Errorf("date is required")
	}

	// check if we can parse date
	time, err := time.Parse(time.RFC3339, form.Value["date"][0])
	if err != nil {
		return nil, fmt.Errorf("error parsing event date: %s", err.Error())
	}

	var description string
	if form.Value["description"] != nil {
		description = form.Value["description"][0]
	}

	// construct the payload
	payload := createEventPayload{
		Name:        form.Value["name"][0],
		Date:        time,
		Description: description,
	}

	return &payload, nil
}
