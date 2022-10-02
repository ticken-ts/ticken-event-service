package dto

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"ticken-event-service/models"
)

type OrganizationDto struct {
	Id    string   `json:"id"`
	Peers []string `json:"peers"`
}

func New(org *models.Organization) OrganizationDto {
	return OrganizationDto{org.OrganizationID, org.Peers}
}

func SendOrganization(ctx *gin.Context, org *models.Organization) {
	dto := New(org)
	res, err := json.Marshal(dto)
	if err != nil {
		ctx.String(500, "error serializing organization")
		ctx.Abort()
		return
	}

	ctx.Data(200, "application/json", res)
}
