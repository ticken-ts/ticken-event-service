package eventController

import (
	"github.com/coreos/go-oidc"
	"github.com/gin-gonic/gin"
	"ticken-event-service/api/dto"
	"ticken-event-service/api/errors"
)

func (controller *EventController) GetEvent(ctx *gin.Context) {
	userID := ctx.MustGet("jwt").(*oidc.IDToken).Subject
	eventId := ctx.Param("eventId")

	event, err := controller.serviceProvider.GetEventManager().GetEvent(eventId, userID)
	if err != nil {
		apiError := errors.GetApiError(err)
		ctx.String(apiError.HttpCode, apiError.Message)
		ctx.Abort()
		return
	}

	dto.SendEvent(ctx, event)
}
