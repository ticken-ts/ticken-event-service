package eventerr

const (
	EventNotFoundErrorCode = iota + 200
	FailedToAddSectionInPVTBC
	EventReadPermissionErrorCode
	StartSaleInPVTBCErrorCode
	StartEventInPVTBCErrorCode
	FinishEventInPVTBCErrorCode
	FailedToStoreEventInPVTBCErrorCode
	FailedToStoreEventInDatabase
)

func GetErrMessage(code uint32) string {
	switch code {
	case EventNotFoundErrorCode:
		return "event not found"
	case FailedToAddSectionInPVTBC:
		return "failed to add section in the private blockchain"
	case EventReadPermissionErrorCode:
		return "organizer is not allowed to read event"
	case StartSaleInPVTBCErrorCode:
		return "failed to set the event on sale in the private blockchain"
	case FailedToStoreEventInPVTBCErrorCode:
		return "an error occurred when trying to store event in private blockchain"
	case FailedToStoreEventInDatabase:
		return "an error occurred when trying to store event in database"
	case StartEventInPVTBCErrorCode:
		return "failed to start event in the private blockchain"
	case FinishEventInPVTBCErrorCode:
		return "failed to finish event in the private blockchain"
	default:
		return "an error has occurred"
	}
}
