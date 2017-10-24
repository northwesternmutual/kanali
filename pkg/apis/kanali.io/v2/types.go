package v2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:noStatus
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ApiKey describes an ApiKey
type ApiKey struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ApiKeySpec `json:"spec"`
}

// APIKeySpec is the spec for an ApiKey
type ApiKeySpec struct {
	Revisions []ApiKeyRevision `json:"revisions"`
}

// ApiKeyRevision is an ApiKey revision
type ApiKeyRevision struct {
	Data     string `json:"data"`
	Status   string `json:"status"`
	LastUsed string `json:"lastUsed"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ApiKeyList is a list of ApiKey resources
type ApiKeyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ApiKey `json:"items"`
}

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MockTarget describes a MockTarget
type MockTarget struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              MockTargetSpec `json:"spec"`
}

// MockTargetSpec is the spec for a MockTarget
type MockTargetSpec struct {
	Routes []Route `json:"routes"`
}

// Route defines the behavior of an http request for a unique route
type Route struct {
	Path       string            `json:"path"`
	StatusCode int               `json:"status"`
	Methods    []string          `json:"methods,omitempty"`
	Headers    map[string]string `json:"headers,omitempty"`
	Body       []byte            `json:"body"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MockTargetList is a list of MockTarget resources
type MockTargetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MockTarget `json:"items"`
}

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ApiProxy describe an ApiProxy
type ApiProxy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ApiProxySpec `json:"spec"`
}

// ApiProxySpec is the spec for an ApiProxy
type ApiProxySpec struct {
	Source  Source   `json:"source"`
	Target  Target   `json:"target"`
	Plugins []Plugin `json:"plugins,omitempty"`
}

// Source represents an incoming request
type Source struct {
	Path        string `json:"path"`
	VirtualHost string `json:"virtualHost,omitempty"`
}

// Target describes a target proxy
type Target struct {
	Path    string  `json:"path,omitempty"`
	Backend Backend `json:"backend,omitempty"`
	SSL     SSL     `json:"ssl,omitempty"`
}

// Mock describes a valid mock response
type Mock struct {
	MockTargetName string `json:"mockTargetName"`
}

// Backend describes an upstream server
type Backend struct {
	Endpoint string  `json:"endpoint,omitempty"`
	Mock     Mock    `json:"mock,omitempty"`
	Service  Service `json:"service,omitempty"`
}

// Service describes a Kubernetes service
type Service struct {
	Name   string  `json:"name,omitempty"`
	Port   int64   `json:"port"`
	Labels []Label `json:"labels,omitempty"`
}

// Label defines a unique label to be matched against service metadata
type Label struct {
	Name   string `json:"name"`
	Header string `json:"header,omitempty"`
	Value  string `json:"value,omitempty"`
}

// SSL describes the anoatomy of a backend TLS connection
type SSL struct {
	SecretName string `json:"secretName"`
}

// Plugin describes a unique plugin to be envoked
type Plugin struct {
	Name    string            `json:"name"`
	Version string            `json:"version,omitempty"`
	Config  map[string]string `json:"config,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ApiProxyList is a list of ApiProxy resources
type ApiProxyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ApiProxy `json:"items"`
}

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ApiKeyBinding describes an ApiKeyBinding
type ApiKeyBinding struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ApiKeyBindingSpec `json:"spec"`
}

// ApiKeyBindingSpec is the spec for an ApiKeyBinding
type ApiKeyBindingSpec struct {
	Keys []Key `json:"keys"`
}

// Key defines a unique key with permissions
type Key struct {
	Name        string `json:"name"`
	Quota       int    `json:"quota,omitempty"`
	Rate        Rate   `json:"rate,omitempty"`
	DefaultRule Rule   `json:"defaultRule,omitempty"`
	Subpaths    []Path `json:"subpaths,omitempty"`
}

// Rate defines rate limit rule
type Rate struct {
	Amount int    `json:"amount,omitempty"`
	Unit   string `json:"unit,omitempty"`
}

// Rule defines a single permission
type Rule struct {
	Global   bool          `json:"global,omitempty"`
	Granular GranularProxy `json:"granular,omitempty"`
}

// Path describes a subpath with unique permissions
type Path struct {
	Path string `json:"path"`
	Rule Rule   `json:"rule,omitempty"`
}

// GranularProxy defines a list of authorized HTTP verbs
type GranularProxy struct {
	Verbs []string `json:"verbs"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ApiKeyBindingList represents a list of ApiKeyBinding resources
type ApiKeyBindingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ApiKeyBinding `json:"items"`
}
