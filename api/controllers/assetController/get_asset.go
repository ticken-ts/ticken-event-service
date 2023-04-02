package assetController

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	"net/http"
	"os"
	"ticken-event-service/utils"
)

func (controller *AssetController) GetAsset(c *gin.Context) {
	assetID, err := uuid.Parse(c.Param("assetID"))
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	asset, err := controller.serviceProvider.GetAssetManager().DownloadAsset(assetID)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
		c.Abort()
		return
	}

	if !isLocalURL(asset.URL) {
		filePath := "/tmp/" + assetID.String()
		err = downloadFile(asset.URL, filePath)
		if err != nil {
			c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
			c.Abort()
			return
		}
		c.File(filePath)
	} else {
		c.File(asset.URL)
	}
}

func isLocalURL(url string) bool {
	return url[0:4] != "http"
}

func downloadFile(url string, filePath string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
