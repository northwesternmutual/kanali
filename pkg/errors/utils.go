package errors

import "fmt"

// ToJSON will turn the passed error into an Error type
// and the JSON data the represents this type.
func ToJSON(err error) (Error, []byte) {
	var e Error
	if _, ok := err.(Error); !ok {
		e = ErrorUnknown
	} else {
		e = err.(Error)
	}
	return e, []byte(fmt.Sprintf(
		`{"status":%d,"message":"%s","code":%d,"details":"%s"}`,
		e.Status,
		e.Message,
		e.Code,
		e.Details,
	))
}
