package validate

import (
	"encoding/json"

	"github.com/northwesternmutual/kanali/pkg/apis/kanali.io/v2"
)

func (v *validation) IsValidMockTarget(data []byte) error {
	mocktarget := new(v2.MockTarget)
	if err := json.Unmarshal(data, mocktarget); err != nil {
		return err
	}
	return v.isValidMockTarget(mocktarget)
}

func (v *validation) isValidMockTarget(mocktarget *v2.MockTarget) error {
	// no dynamic validation needed at this time
	return nil
}

func (v *validation) IsValidMockTargetList(data []byte) error {
	list := new(v2.MockTargetList)
	if err := json.Unmarshal(data, list); err != nil {
		return err
	}
	return v.isValidMockTargetList(list)
}

func (v *validation) isValidMockTargetList(list *v2.MockTargetList) error {
	for _, mocktarget := range list.Items {
		if err := v.isValidMockTarget(&mocktarget); err != nil {
			return err
		}
	}
	return nil
}
