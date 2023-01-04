package exception

import "fmt"

func FromError(err error, msg string) error {
	return fmt.Errorf("[Ticken Error - %s] - [%s]", msg, err.Error())
}

func WithMessage(format string, a ...any) error {
	msg := fmt.Sprintf(format, a)
	return fmt.Errorf("[Ticken Error - %s]", msg)
}
