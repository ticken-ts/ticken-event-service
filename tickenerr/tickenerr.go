package tickenerr

import (
	"fmt"
	"strconv"
	"ticken-event-service/tickenerr/asseterr"
	"ticken-event-service/tickenerr/commonerr"
	"ticken-event-service/tickenerr/eventerr"
	"ticken-event-service/tickenerr/organizationerr"
	"ticken-event-service/tickenerr/organizererr"
)

type TickenError struct {
	Message       string
	Code          uint32
	UnderlyingErr error
}

func New(errCode uint32) TickenError {
	return FromErrorWithMessage(errCode, nil, "")
}

func NewWithMessage(errCode uint32, msg string) TickenError {
	return FromErrorWithMessage(errCode, nil, msg)
}

func FromError(errCode uint32, underlyingError error) TickenError {
	return FromErrorWithMessage(errCode, underlyingError, "")
}

func FromErrorWithMessage(errCode uint32, underlyingError error, extraMsg string) TickenError {
	var message string

	if between(errCode, 0, 99) {
		message = commonerr.GetErrMessage(errCode)
	}
	if between(errCode, 100, 199) {
		message = organizererr.GetErrMessage(errCode)
	}
	if between(errCode, 200, 299) {
		message = eventerr.GetErrMessage(errCode)
	}
	if between(errCode, 400, 499) {
		message = organizationerr.GetErrMessage(errCode)
	}
	if between(errCode, 500, 599) {
		message = asseterr.GetErrMessage(errCode)
	}

	if len(extraMsg) > 0 {
		message = fmt.Sprintf("%s (%s)", message, extraMsg)
	}

	return TickenError{
		Message:       message,
		Code:          errCode,
		UnderlyingErr: underlyingError,
	}
}

func between(code, min, max uint32) bool {
	return code >= min && code <= max
}

func (ticketErr TickenError) Error() string {
	return strconv.Itoa(int(ticketErr.Code))
}
