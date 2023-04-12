package eventController

import (
	"net/http"
	"ticken-event-service/api/mappers"
	"ticken-event-service/api/res"
	"ticken-event-service/security/jwt"

	"github.com/gin-gonic/gin"
)

func (controller *EventController) GetMyOrganizations(c *gin.Context) {
	organizerID := c.MustGet("jwt").(*jwt.Token).Subject

	organizations, err := controller.ServiceProvider.GetOrganizationManager().GetOrganizationsByOrganizer(organizerID)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, res.Success{
		Message: "organizations retrieved successfully",
		Data:    mappers.MapOrganizationListToDTO(organizations),
	})
}
