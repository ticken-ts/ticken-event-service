package organizationController

import (
	"github.com/coreos/go-oidc"
	"github.com/gin-gonic/gin"
	"ticken-event-service/api/dto"
	"ticken-event-service/api/errors"
)

func (controller *OrganizationController) GetMyOrganization(ctx *gin.Context) {
	userID := ctx.MustGet("jwt").(*oidc.IDToken).Subject
	org, err := controller.serviceProvider.GetOrgManager().GetUserOrganization(userID)

	if err != nil {
		apiError := errors.GetApiError(err)
		ctx.String(apiError.HttpCode, apiError.Message)
		ctx.Abort()
	}

	dto.SendOrganization(ctx, org)
}
