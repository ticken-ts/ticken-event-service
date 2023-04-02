package asseterr

const (
	AssetNotFoundErrorCode = iota
	FailedToStoreAssetInDatabase
	FailedToUploadAsset
	FailedToDownloadAsset
)

func GetErrMessage(code uint32) string {
	switch code {
	case AssetNotFoundErrorCode:
		return "asset not found"
	case FailedToStoreAssetInDatabase:
		return "an error occurred while storing asset info in database"
	case FailedToUploadAsset:
		return "an error occurred while uploading asset to bucket"
	case FailedToDownloadAsset:
		return "an error occurred while downloading asset to bucket"
	default:
		return "an error has occurred"
	}
}
