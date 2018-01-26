package middleware

import (
  "time"
	"testing"
  "net/http/httptest"

	"github.com/stretchr/testify/assert"
)

func TestGetRequDuration(t *testing.T) {
	t0, t1 := time.Date(2018, time.January, 1, 2, 3, 4, 1000000, time.Local), time.Date(2018, time.January, 1, 2, 3, 4, 256000000, time.Local)
	assert.Equal(t, float64(255), getReqDuration(t0, t1))
}

func TestGetRemoteAddr(t *testing.T) {
	assert.Equal(t, "", getRemoteAddr(""))
	assert.Equal(t, "1.2.3.4", getRemoteAddr("1.2.3.4"))
	assert.Equal(t, "1.2.3.4", getRemoteAddr("1.2.3.4:8080"))
	assert.Equal(t, "foo", getRemoteAddr("foo"))
	assert.Equal(t, "[::]:8080", getRemoteAddr("[::]:8080"))
}