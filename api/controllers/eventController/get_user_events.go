package eventController

import (
	"github.com/gin-gonic/gin"
	"ticken-event-service/api/dto"
	"ticken-event-service/api/errors"
)

func (controller *EventController) GetUserEvents(ctx *gin.Context) {
	userId := ctx.GetString("userId")

	events, err := controller.serviceProvider.GetEventManager().GetUserEvents(userId)
	if err != nil {
		apiError := errors.GetApiError(err)
		ctx.String(apiError.HttpCode, apiError.Message)
		ctx.Abort()
		return
	}

	dto.SendEvents(ctx, events)
}
