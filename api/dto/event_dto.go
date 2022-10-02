package dto

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"ticken-event-service/models"
)

func SendEvent(ctx *gin.Context, event *models.Event) {

	res, err := json.Marshal(event)
	if err != nil {
		ctx.String(500, "error serializing event")
		ctx.Abort()
		return
	}

	ctx.Data(200, "application/json", res)

}

func SendEvents(ctx *gin.Context, events []*models.Event) {
	res, err := json.Marshal(events)
	if err != nil {
		ctx.String(500, "error serializing event")
		ctx.Abort()
		return
	}

	ctx.Data(200, "application/json", res)

}
