package errors

type ApiError struct {
	Message  string
	HttpCode int
}

const (
	UserOrgNotFound  = "USER_ORG_NOT_FOUND"
	EventNotFound    = "EVENT_NOT_FOUND"
	OrgEventMismatch = "ORG_EVENT_MISMATCH"
	InvalidToken     = "INVALID_TOKEN"
)

func newApiError(message string, httpCode int) *ApiError {
	err := new(ApiError)
	err.Message = message
	err.HttpCode = httpCode
	return err
}

func GetApiError(err error) *ApiError {
	var errors = map[string]*ApiError{
		UserOrgNotFound:  newApiError("user organization not found", 404),
		EventNotFound:    newApiError("event not found", 404),
		OrgEventMismatch: newApiError("user organization does not own event", 401),
		InvalidToken:     newApiError("invalid token", 401),
	}

	return errors[err.Error()]
}
