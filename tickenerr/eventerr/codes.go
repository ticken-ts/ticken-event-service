package eventerr

const (
	EventNotFoundErrorCode = iota + 200
	CreateEventErrorCode
	FailedToAddSectionInPVTBC
	EventReadPermissionErrorCode
	SetTicketOnSaleInPVTBCErrorCode
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
	case SetTicketOnSaleInPVTBCErrorCode:
		return "failed to set the event on sale in the private blockchain"
	case FailedToStoreEventInPVTBCErrorCode:
		return "an error occurred when trying to store event in private blockchain"
	case FailedToStoreEventInDatabase:
		return "an error occurred when trying to store event in database"
	default:
		return "an error has occurred"
	}
}
