package organizationerr

const (
	OrganizationNotFoundErrorCode = iota + 400
	RegisterValidatorErrorCode
	EstablishPVTBCConnectionErrorCode
)

func GetErrMessage(code uint32) string {
	switch code {
	case OrganizationNotFoundErrorCode:
		return "organization not found"
	case RegisterValidatorErrorCode:
		return "failed to register validator for organization"
	case EstablishPVTBCConnectionErrorCode:
		return "failed to establish private blockchain connection error"
	default:
		return "an error has occurred"
	}
}
