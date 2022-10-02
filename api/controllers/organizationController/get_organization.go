package organizationController

import (
	"github.com/gin-gonic/gin"
	"ticken-event-service/api/dto"
	"ticken-event-service/api/errors"
)

func (controller *OrganizationController) GetMyOrganization(ctx *gin.Context) {
	userId := ctx.GetString("userId")
	org, err := controller.serviceProvider.GetOrgManager().GetUserOrganization(userId)

	if err != nil {
		apiError := errors.GetApiError(err)
		ctx.String(apiError.HttpCode, apiError.Message)
		ctx.Abort()
	}

	dto.SendOrganization(ctx, org)
}
