package eventController

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
)

func (controller *EventController) GetEvent(ctx *gin.Context) {
	eventId := ctx.Param("eventId")
	userId := ctx.GetString("userId")

	event, err := controller.serviceProvider.GetEventManager().GetEvent(eventId, userId)
	if err != nil {
		ctx.String(404, err.Error())
		ctx.Abort()
		return
	}

	res, err := json.Marshal(event)
	if err != nil {
		ctx.String(500, "error serializing event")
		ctx.Abort()
		return
	}

	ctx.Data(200, "application/json", res)
}
