package eventController

import (
	"github.com/gin-gonic/gin"
	"ticken-event-service/api/dto"
	"ticken-event-service/api/errors"
)

func (controller *EventController) GetEvent(ctx *gin.Context) {
	eventId := ctx.Param("eventId")
	userId := ctx.GetString("userId")

	event, err := controller.serviceProvider.GetEventManager().GetEvent(eventId, userId)
	if err != nil {
		apiError := errors.GetApiError(err)
		ctx.String(apiError.HttpCode, apiError.Message)
		ctx.Abort()
		return
	}

	dto.SendEvent(ctx, event)
}
