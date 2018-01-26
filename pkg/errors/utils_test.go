package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToJSON(t *testing.T) {
	typed := Error{404, "message", 01, "details"}
	untyped := errors.New("untyped")

	e, d := ToJSON(typed)
	assert.Equal(t, e, typed)
	assert.Equal(t, d, []byte(`{"status":404,"message":"message","code":01,"details":"details"}`))

	e, d = ToJSON(untyped)
	assert.Equal(t, e, ErrorUnknown)
	assert.Equal(t, d, []byte(`{"status":500,"message":"An unknown error occured.","code":01,"details":"More details coming soon!"}`))
}
