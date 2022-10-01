package eventController

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"ticken-event-service/models"
)

func (controller *EventController) GetEvent(ctx *gin.Context) {
	res, err := json.Marshal(models.Event{})
	if err != nil {
		panic("error returning event")
	}
	ctx.Data(200, "application/json", res)
}
