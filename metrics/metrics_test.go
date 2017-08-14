package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	f := &Metrics{}
	f.Add(Metric{"nameOne", "valueOne", false})
	f.Add(Metric{"nameTwo", "valueTwo", true}, Metric{"nameThree", "valueThree", false})
	assert.Equal(t, len(*f), 3)
	assert.False(t, (*f)[0].Index)
	assert.True(t, (*f)[1].Index)
	assert.False(t, (*f)[2].Index)
	assert.Equal(t, (*f)[0].Name, "nameOne")
	assert.Equal(t, (*f)[1].Name, "nameTwo")
	assert.Equal(t, (*f)[2].Name, "nameThree")
	assert.Equal(t, (*f)[0].Value, "valueOne")
	assert.Equal(t, (*f)[1].Value, "valueTwo")
	assert.Equal(t, (*f)[2].Value, "valueThree")
}
