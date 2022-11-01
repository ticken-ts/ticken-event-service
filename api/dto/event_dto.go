package dto

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"ticken-event-service/models"
)

type EventDTO struct {
	EventID string `json:"event_id"`
	Name    string `json:"name"`
	Date    string `json:"date"`
	OnChain bool   `json:"on_chain"`
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
