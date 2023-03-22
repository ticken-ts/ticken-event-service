package validatorController

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"ticken-event-service/api/mappers"
	"ticken-event-service/security/jwt"
	"ticken-event-service/utils"
)

type validatorPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// transferencia
// scannear tickets -> app tickets

func (controller *ValidatorController) RegisterValidator(c *gin.Context) {
	organizerID := c.MustGet("jwt").(*jwt.Token).Subject

	var payload validatorPayload

	organizationID, err := uuid.Parse(c.Param("organizationID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
		c.Abort()
		return
	}

	if err = c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
		c.Abort()
		return
	}

	validatorManager := controller.serviceProvider.GetValidatorManager()

	newValidator, err := validatorManager.RegisterValidator(
		organizerID,
		organizationID,
		payload.Username,
		payload.Password,
		payload.Email,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, utils.HttpResponse{
		Message: "Validator created",
		Data:    mappers.MapValidatorToValidatorDTO(newValidator),
	})
}
