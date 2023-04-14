package eventController

import (
	"mime/multipart"
	"net/http"
	"ticken-event-service/api/mappers"
	"ticken-event-service/api/res"
	"ticken-event-service/security/jwt"
	"ticken-event-service/utils/file"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type createEventPayload struct {
	Name        string                `form:"name" binding:"required"`
	Date        time.Time             `form:"date" binding:"required"`
	Description string                `form:"description" binding:"required"`
	PosterFile  *multipart.FileHeader `form:"poster"`
}

func (controller *EventController) CreateEvent(c *gin.Context) {
	userID := c.MustGet("jwt").(*jwt.Token).Subject

	var payload createEventPayload

	organizationID, err := uuid.Parse(c.Param("organizationID"))
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	// validate the form
	if err := c.Bind(&payload); err != nil {
		c.Abort()
		return
	}

	// the only thing that we are going to validate
	// is the that we can bind the request to the struct
	// No further validation are going to be added, so all
	// validations are going to be performed on chain

	var file *file.File
	if payload.PosterFile != nil {
		file, err = controller.ReadFile(payload.PosterFile)
		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}
	}

	event, err := controller.ServiceProvider.GetEventManager().CreateEvent(
		userID,
		organizationID,
		payload.Name,
		payload.Date,
		payload.Description,
		file,
	)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, res.Success{
		Message: "event created successfully",
		Data:    mappers.MapEventToEventDTO(event),
	})
}
