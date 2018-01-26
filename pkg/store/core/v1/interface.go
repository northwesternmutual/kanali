package v1

import (
	"k8s.io/client-go/informers/core/v1"
)

var v1Interface v1.Interface

func SetGlobalInterface(i v1.Interface) {
	v1Interface = i
}

func Interface() v1.Interface {
	return v1Interface
}
