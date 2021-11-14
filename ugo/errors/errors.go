package errors

import (
	"fmt"
)

type Status int

type Error struct {
	Pos     string
	Message string
}

func New(pos string, a ...interface{}) error {
	return &Error{Pos: pos, Message: fmt.Sprint(a...)}
}
func Newf(pos string, format string, a ...interface{}) error {
	return &Error{Pos: pos, Message: fmt.Sprintf(format, a...)}
}

func (p *Error) Error() string {
	if p.Pos != "" {
		return fmt.Sprintf("%s: %s", p.Pos, p.Message)
	} else {
		return p.Message
	}
}
