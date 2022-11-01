package eventController

import (
	"github.com/coreos/go-oidc"
	"github.com/gin-gonic/gin"
	"net/http"
	"ticken-event-service/api/mappers"
	"ticken-event-service/utils"
)

func (controller *EventController) GetUserEvents(ctx *gin.Context) {
	userID := ctx.MustGet("jwt").(*oidc.IDToken).Subject

	events, err := controller.serviceProvider.GetEventManager().GetOrganizationEvents(userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
		ctx.Abort()
		return
	}

	eventsDTO := mappers.MapEventListToEventListDTO(events)

	ctx.JSON(http.StatusOK, utils.HttpResponse{Data: eventsDTO})
}
