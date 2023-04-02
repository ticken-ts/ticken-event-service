package organizererr

const (
	OrganizerNotFoundErrorCode = iota + 100
	OrganizerNotBelongsToOrganization
)

func GetErrMessage(code uint32) string {
	switch code {
	case OrganizerNotFoundErrorCode:
		return "organizer not found"
	case OrganizerNotBelongsToOrganization:
		return "organizer doest not belongs to organization"
	default:
		return "an error has occurred"
	}
}
