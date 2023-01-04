package organizationController

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"ticken-event-service/api/security"
	"ticken-event-service/utils"
)

func (controller *OrganizationController) GetOrganizationCrypto(c *gin.Context) {
	organizationID := c.Param("organizationID")
	jwt := c.MustGet("jwt").(*security.JWT)

	organizationManager := controller.serviceProvider.GetOrganizationManager()

	zip, err := organizationManager.GetOrganizationCryptoZipped(jwt.Subject, organizationID)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
		c.Abort()
		return
	}

	c.Writer.Header().Set("Content-Type", "application/zip")
	c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.zip\"", "crypto"))
	_, _ = c.Writer.Write(zip)
}
