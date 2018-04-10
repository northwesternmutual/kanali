package validate

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/northwesternmutual/kanali/pkg/apis/kanali.io/v2"
	"github.com/northwesternmutual/kanali/test/builder"
)

func TestIsValidMockTarget(t *testing.T) {
	assert.Error(t, (&validation{}).IsValidMockTarget(nil))

	mocktarget, _ := json.Marshal(builder.NewMockTarget("foo", "bar").NewOrDie())
	assert.Nil(t, (&validation{}).IsValidMockTarget(mocktarget))
}

func TestIsValidMockTargetList(t *testing.T) {
	assert.Error(t, (&validation{}).IsValidMockTargetList(nil))

	list, _ := json.Marshal(&v2.MockTargetList{
		Items: []v2.MockTarget{
			*builder.NewMockTarget("foo", "bar").NewOrDie(),
		},
	})
	assert.Nil(t, (&validation{}).IsValidMockTargetList(list))
}
