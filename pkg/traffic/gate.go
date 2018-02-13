package traffic

import (
	"github.com/northwesternmutual/kanali/pkg/feature"
)

const (
	featureName = "RateLimiter"
)

func (ctlr *Controller) FeatureIsEnabled() bool {
	return feature.Contains(featureName)
}
