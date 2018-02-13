package feature

import "strings"

// Feature is an interface that describes if a feature
// is enabled or not. It also specifies the name of the
// feature so that it can be assoicated with an application's
// configuration.
type Feature interface {
	FeatureIsEnabled() bool
}

var enabledFeatures = []string{}

func SetEnabled(l []string) {
	enabledFeatures = l
}

func Contains(featureName string) bool {
	for _, f := range enabledFeatures {
		if strings.ToUpper(f) == strings.ToUpper(featureName) {
			return true
		}
	}
	return false
}
