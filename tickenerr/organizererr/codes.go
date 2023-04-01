package organizererr

const (
	OrganizerNotFoundErrorCode = iota + 100
)

func GetErrMessage(code uint32) string {
	switch code {
	case OrganizerNotFoundErrorCode:
		return "organizer not found"
	default:
		return "an error has occurred"
	}
}
