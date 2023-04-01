package eventController

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"ticken-event-service/api/mappers"
	"ticken-event-service/api/res"
	"ticken-event-service/security/jwt"
)

func (controller *EventController) GetOrganizationEvents(c *gin.Context) {
	organizerID := c.MustGet("jwt").(*jwt.Token).Subject

	organizationID, err := uuid.Parse(c.Param("organizationID"))
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	events, err := controller.serviceProvider.GetEventManager().GetOrganizationEvents(
		organizerID,
		organizationID,
	)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, res.Success{
		Message: fmt.Sprintf("%d events found", len(events)),
		Data:    mappers.MapEventListToDTO(events),
	})
}
