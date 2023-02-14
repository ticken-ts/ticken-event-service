package publicController

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"ticken-event-service/api/mappers"
	"ticken-event-service/utils"
)

func (controller *PublicController) GetAvailableEvents(c *gin.Context) {
	manager := controller.serviceProvider.GetEventManager()
	events, err := manager.GetAvailableEvents()
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
		c.Abort()
		return
	}

	eventsDTO := mappers.MapEventListToDTO(events)

	c.JSON(http.StatusOK, utils.HttpResponse{Data: eventsDTO})

}
