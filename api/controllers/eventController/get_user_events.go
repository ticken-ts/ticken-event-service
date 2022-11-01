package eventController

import (
	"github.com/coreos/go-oidc"
	"github.com/gin-gonic/gin"
	"ticken-event-service/api/dto"
	"ticken-event-service/api/errors"
)

func (controller *EventController) GetUserEvents(ctx *gin.Context) {
	userID := ctx.MustGet("jwt").(*oidc.IDToken).Subject

	events, err := controller.serviceProvider.GetEventManager().GetOrganizationEvents(userID)
	if err != nil {
		apiError := errors.GetApiError(err)
		ctx.String(apiError.HttpCode, apiError.Message)
		ctx.Abort()
		return
	}

	dto.SendEvents(ctx, events)
}
