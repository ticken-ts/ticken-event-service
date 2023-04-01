package eventerr

const (
	EventNotFoundErrorCode = iota + 200
	FailedToAddSectionInPVTBC
	EventReadPermissionErrorCode
	SetTicketOnSaleInPVTBCErrorCode
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
		return "failed to set the evebt on sale in the private blockchain"
	default:
		return "an error has occurred"
	}
}
