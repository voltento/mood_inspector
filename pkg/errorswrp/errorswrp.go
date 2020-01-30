package errorswrp

import (
	"errors"
	"fmt"
)

func Wrap(err error, msg string) error {
	newMsg := fmt.Sprintf("%v: %v", err.Error(), msg)
	return errors.New(newMsg)
}
