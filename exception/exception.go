package exception

import "fmt"

func WithMessage(err error, msg string) error {
	return fmt.Errorf("[Ticken Error - %s] - [%s]", msg, err.Error())
}
