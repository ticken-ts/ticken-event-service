package organizerController

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"ticken-event-service/api/mappers"
	"ticken-event-service/utils"
)

func (controller *OrganizerController) RegisterOrganizer(c *gin.Context) {
	jwt := c.MustGet("jwt")
	userID := utils.GetUserIDFromJWT(jwt)

	claims, err := utils.ParseJWTClaims(jwt)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: "could not get claims from jwt: " + err.Error()})
		c.Abort()
		return
	}

	organizerManager := controller.serviceProvider.GetOrganizationManager()

	//
	organizer, err := organizerManager.RegisterOrganizer(userID, claims.Email, claims.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
		c.Abort()
		return
	}

	organizerDTO := mappers.MapOrganizerToOrganizerDTO(organizer)

	c.JSON(http.StatusCreated, utils.HttpResponse{Data: organizerDTO})
}
