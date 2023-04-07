package publicController

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"ticken-event-service/api/mappers"
	"ticken-event-service/api/res"
	"time"
)

func (controller *PublicController) GetEventsOnSale(c *gin.Context) {
	eventName := c.Param("name")

	fromDate, err := parseDateOrZero(c.Param("fromDate"))
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	toDate, err := parseDateOrZero(c.Param("toDate"))
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	events, err := controller.serviceProvider.GetEventManager().GetEventsOnSale(
		eventName,
		fromDate,
		toDate,
	)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, res.Success{
		Message: fmt.Sprintf("%d events found", len(events)),
		Data:    mappers.MapEventListToDTO(events),
	})
}

func parseDateOrZero(dateStr string) (time.Time, error) {
	var date time.Time

	if len(dateStr) > 0 {
		date, err := time.Parse(time.DateOnly, dateStr)
		if err != nil {
			return date, fmt.Errorf("invalid date format: %s", err.Error())
		}
	}

	return date, nil
}
