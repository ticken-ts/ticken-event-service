package assetController

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (controller *AssetController) GetAsset(c *gin.Context) {
	assetID, err := uuid.Parse(c.Param("assetID"))
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	assetURL, err := controller.serviceProvider.GetAssetManager().GetAssetURL(assetID)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	// todo -> check if works with external assets
	// stored in aws or cdn
	c.File(assetURL)
}
