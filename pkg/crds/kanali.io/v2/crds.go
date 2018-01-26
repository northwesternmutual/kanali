// Copyright (c) 2017 Northwestern Mutual.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package v2

import (
	"fmt"

  metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"


  kanaliio "github.com/northwesternmutual/kanali/pkg/apis/kanali.io"
)

const (
	version         = "v2"
)

var (
	rfc3339Regex       = `^([0-9]+)-(0[1-9]|1[012])-(0[1-9]|[12][0-9]|3[01])[Tt]([01][0-9]|2[0-3]):([0-5][0-9]):([0-5][0-9]|60)(\.[0-9]+)?(([Zz])|([\+|\-]([01][0-9]|2[0-3]):[0-5][0-9]))$`
	encryptedDataRegex = `[0-9a-zA-Z]+`
	httpPathRegex      = `^\/.*`
	// https://stackoverflow.com/questions/1418423/the-hostname-regex
	virtualHostRegex = `^(?=.{1,255}$)[0-9A-Za-z](?:(?:[0-9A-Za-z]|-){0,61}[0-9A-Za-z])?(?:\.[0-9A-Za-z](?:(?:[0-9A-Za-z]|-){0,61}[0-9A-Za-z])?)*\.?$`
	// https://stackoverflow.com/questions/3809401/what-is-a-good-regular-expression-to-match-a-url
	httpURLRegex = `https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,6}\b([-a-zA-Z0-9@:%_\+.~#?&//=]*)`
	// https://github.com/kubernetes/apimachinery/blob/master/pkg/util/validation/validation.go
	k8sPrefixFmt          = `([a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*\/?)?`
	k8sNameFmt            = `([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9]`
	k8sNameRegex          = `^` + k8sNameFmt + `$`
	k8sQualifiedNameRegex = `^` + k8sPrefixFmt + k8sNameFmt + `$`
	httpHeaderNameRegex   = `^[0-9a-zA-Z](.*)[0-9a-zA-Z]$`
	semanticVersionRegex  = `^v?(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(-(0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(\.(0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*)?(\+[0-9a-zA-Z-]+(\.[0-9a-zA-Z-]+)*)?$`
)

var ApiKeyCRD = &apiextensionsv1beta1.CustomResourceDefinition{
	ObjectMeta: metav1.ObjectMeta{
		Name: fmt.Sprintf("apikeys.%s", kanaliio.GroupName),
	},
	Spec: apiextensionsv1beta1.CustomResourceDefinitionSpec{
		Group:   kanaliio.GroupName,
		Version: version,
		Names: apiextensionsv1beta1.CustomResourceDefinitionNames{
			Plural:   "apikeys",
			Singular: "apikey",
			ShortNames: []string{
				"ak",
			},
			Kind:     "ApiKey",
			ListKind: "ApiKeyList",
		},
		Scope: apiextensionsv1beta1.ClusterScoped,
		Validation: &apiextensionsv1beta1.CustomResourceValidation{
			OpenAPIV3Schema: &apiextensionsv1beta1.JSONSchemaProps{
				Description: "top level field wrapping for the ApiKey resource",
				Required: []string{
					"spec",
				},
				Properties: map[string]apiextensionsv1beta1.JSONSchemaProps{
					"spec": {
						Description: "ApiKey resource body",
						Required: []string{
							"revisions",
						},
						Properties: map[string]apiextensionsv1beta1.JSONSchemaProps{
							"revisions": {
								Description: "represents the list of active and inactive API key revisions",
								Type:        "array",
								MinLength:   int64Ptr(1),
								Items: &apiextensionsv1beta1.JSONSchemaPropsOrArray{
									Schema: &apiextensionsv1beta1.JSONSchemaProps{
										Description: "a represetation of an API key revision",
										Type:        "object",
										Required: []string{
											"data",
											"status",
										},
										Properties: map[string]apiextensionsv1beta1.JSONSchemaProps{
											"data": {
												Description: "rsa encrypted API key data",
												Type:        "string",
												MinLength:   int64Ptr(1),
												Pattern:     encryptedDataRegex,
											},
											"status": {
												Description: "status of API key",
												Type:        "string",
												Enum: []apiextensionsv1beta1.JSON{
													{
														Raw: []byte(`"Active"`),
													},
													{
														Raw: []byte(`"Inactive"`),
													},
												},
											},
											"lastUsed": {
												Description: "RFC3339 timestamp that this API key revision was last used in an attempted request",
												Type:        "string",
												// regex from https://gist.github.com/marcelotmelo/b67f58a08bee6c2468f8
												Pattern: rfc3339Regex,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	},
}

var ApiKeyBindingCRD = &apiextensionsv1beta1.CustomResourceDefinition{
	ObjectMeta: metav1.ObjectMeta{
		Name: fmt.Sprintf("apikeybindings.%s", kanaliio.GroupName),
	},
	Spec: apiextensionsv1beta1.CustomResourceDefinitionSpec{
		Group:   kanaliio.GroupName,
		Version: version,
		Names: apiextensionsv1beta1.CustomResourceDefinitionNames{
			Plural:   "apikeybindings",
			Singular: "apikeybinding",
			ShortNames: []string{
				"akb",
			},
			Kind:     "ApiKeyBinding",
			ListKind: "ApiKeyBindingList",
		},
		Scope: apiextensionsv1beta1.NamespaceScoped,
		Validation: &apiextensionsv1beta1.CustomResourceValidation{
			OpenAPIV3Schema: &apiextensionsv1beta1.JSONSchemaProps{
				Description: "top level field wrapping for the ApiKeyBinding resource",
				Required: []string{
					"spec",
				},
				Properties: map[string]apiextensionsv1beta1.JSONSchemaProps{
					"spec": {
						Description: "ApiKeyBinding resource body",
						Required: []string{
							"keys",
						},
						Properties: map[string]apiextensionsv1beta1.JSONSchemaProps{
							"keys": {
								Description: "list of ApiKey resources granted permissions",
								Type:        "array",
								MinLength:   int64Ptr(1),
								Items: &apiextensionsv1beta1.JSONSchemaPropsOrArray{
									Schema: &apiextensionsv1beta1.JSONSchemaProps{
										Description: "ApiKey permissions",
										Type:        "object",
										Required: []string{
											"name",
										},
										Properties: map[string]apiextensionsv1beta1.JSONSchemaProps{
											"name": {
												Description: "valid qualified Kubernetes name",
												Type:        "string",
												MinLength:   int64Ptr(1),
												MaxLength:   int64Ptr(63),
												Pattern:     k8sQualifiedNameRegex,
											},
											"defaultRule": {
												Description: "defines http method that an ApiKey has access to",
												Type:        "object",
												OneOf: []apiextensionsv1beta1.JSONSchemaProps{
													{
														Description: "global access",
														Type:        "object",
														Required: []string{
															"global",
														},
														Properties: map[string]apiextensionsv1beta1.JSONSchemaProps{
															"global": {
																Description: "does ApiKey have access to all http methods",
																Type:        "boolean",
															},
														},
													},
													{
														Description: "fine grained http verb access",
														Type:        "object",
														Required: []string{
															"granular",
														},
														Properties: map[string]apiextensionsv1beta1.JSONSchemaProps{
															"granular": {
																Description: "fine grained http verb access",
																Type:        "object",
																Required: []string{
																	"verbs",
																},
																Properties: map[string]apiextensionsv1beta1.JSONSchemaProps{
																	"verbs": {
																		Description: "list of http methods that ApiKey has access to",
																		Type:        "array",
																		MinLength:   int64Ptr(1),
																		Items: &apiextensionsv1beta1.JSONSchemaPropsOrArray{
																			Schema: &apiextensionsv1beta1.JSONSchemaProps{
																				Description: "valid htp methods",
																				Type:        "string",
																				Enum: []apiextensionsv1beta1.JSON{
																					{
																						Raw: []byte(`"GET"`),
																					},
																					{
																						Raw: []byte(`"HEAD"`),
																					},
																					{
																						Raw: []byte(`"POST"`),
																					},
																					{
																						Raw: []byte(`"PUT"`),
																					},
																					{
																						Raw: []byte(`"PATCH"`),
																					},
																					{
																						Raw: []byte(`"DELETE"`),
																					},
																					{
																						Raw: []byte(`"CONNECT"`),
																					},
																					{
																						Raw: []byte(`"OPTIONS"`),
																					},
																					{
																						Raw: []byte(`"TRACE"`),
																					},
																				},
																			},
																		},
																	},
																},
															},
														},
													},
												},
											},
											"subpaths": {
												Description: "list of subpaths defining fine grained permissions",
												Type:        "array",
												MinLength:   int64Ptr(1),
												Items: &apiextensionsv1beta1.JSONSchemaPropsOrArray{
													Schema: &apiextensionsv1beta1.JSONSchemaProps{
														Description: "unique subpath",
														Required: []string{
															"path",
															"rule",
														},
														Properties: map[string]apiextensionsv1beta1.JSONSchemaProps{
															"path": {
																Description: "http path",
																Type:        "string",
																Pattern:     httpPathRegex,
																MinLength:   int64Ptr(1),
																Default: &apiextensionsv1beta1.JSON{
																	Raw: []byte(""),
																},
															},
															"rule": {
																Description: "defines http method that an ApiKey has access to",
																Type:        "object",
																OneOf: []apiextensionsv1beta1.JSONSchemaProps{
																	{
																		Description: "global access",
																		Type:        "object",
																		Required: []string{
																			"global",
																		},
																		Properties: map[string]apiextensionsv1beta1.JSONSchemaProps{
																			"global": {
																				Description: "does ApiKey have access to all http methods",
																				Type:        "boolean",
																			},
																		},
																	},
																	{
																		Description: "fine grained http verb access",
																		Type:        "object",
																		Required: []string{
																			"granular",
																		},
																		Properties: map[string]apiextensionsv1beta1.JSONSchemaProps{
																			"granular": {
																				Description: "fine grained http verb access",
																				Type:        "object",
																				Required: []string{
																					"verbs",
																				},
																				Properties: map[string]apiextensionsv1beta1.JSONSchemaProps{
																					"verbs": {
																						Description: "list of http methods that ApiKey has access to",
																						Type:        "array",
																						MinLength:   int64Ptr(1),
																						Items: &apiextensionsv1beta1.JSONSchemaPropsOrArray{
																							Schema: &apiextensionsv1beta1.JSONSchemaProps{
																								Description: "valid htp methods",
																								Type:        "string",
																								Enum: []apiextensionsv1beta1.JSON{
																									{
																										Raw: []byte(`"GET"`),
																									},
																									{
																										Raw: []byte(`"HEAD"`),
																									},
																									{
																										Raw: []byte(`"POST"`),
																									},
																									{
																										Raw: []byte(`"PUT"`),
																									},
																									{
																										Raw: []byte(`"PATCH"`),
																									},
																									{
																										Raw: []byte(`"DELETE"`),
																									},
																									{
																										Raw: []byte(`"CONNECT"`),
																									},
																									{
																										Raw: []byte(`"OPTIONS"`),
																									},
																									{
																										Raw: []byte(`"TRACE"`),
																									},
																								},
																							},
																						},
																					},
																				},
																			},
																		},
																	},
																},
															},
														},
													},
												},
											},
											"rate": {
												Description: "number of requests an ApiKey can make over an interval",
												Type:        "object",
												Properties: map[string]apiextensionsv1beta1.JSONSchemaProps{
													"amount": {
														Type:       "integer",
														Minimum:    float64Ptr(1),
														MultipleOf: float64Ptr(1),
													},
													"unit": {
														Description: "valid intervals",
														Type:        "string",
														Enum: []apiextensionsv1beta1.JSON{
															{
																Raw: []byte(`"SECOND"`),
															},
															{
																Raw: []byte(`"MINUTE"`),
															},
															{
																Raw: []byte(`"HOUR"`),
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	},
}

var ApiProxyCRD = &apiextensionsv1beta1.CustomResourceDefinition{
	ObjectMeta: metav1.ObjectMeta{
		Name: fmt.Sprintf("apiproxies.%s", kanaliio.GroupName),
	},
	Spec: apiextensionsv1beta1.CustomResourceDefinitionSpec{
		Group:   kanaliio.GroupName,
		Version: version,
		Names: apiextensionsv1beta1.CustomResourceDefinitionNames{
			Plural:   "apiproxies",
			Singular: "apiproxy",
			ShortNames: []string{
				"ap",
			},
			Kind:     "ApiProxy",
			ListKind: "ApiProxyList",
		},
		Scope: apiextensionsv1beta1.NamespaceScoped,
		Validation: &apiextensionsv1beta1.CustomResourceValidation{
			OpenAPIV3Schema: &apiextensionsv1beta1.JSONSchemaProps{
				Description: "top level field wrapping for the ApiProxy resource",
				Required: []string{
					"spec",
				},
				Properties: map[string]apiextensionsv1beta1.JSONSchemaProps{
					"spec": {
						Description: "ApiProxy resource body",
						Required: []string{
							"source",
							"target",
						},
						Properties: map[string]apiextensionsv1beta1.JSONSchemaProps{
							"source": {
								Description: "unique incoming http path and virtual host combination",
								Required: []string{
									"path",
								},
								Properties: map[string]apiextensionsv1beta1.JSONSchemaProps{
									"path": {
										Description: "http path",
										Type:        "string",
										Pattern:     httpPathRegex,
										MinLength:   int64Ptr(1),
										Default: &apiextensionsv1beta1.JSON{
											Raw: []byte(""),
										},
									},
									"virtualHost": {
										Description: "http hostname",
										Type:        "string",
										Pattern:     virtualHostRegex,
									},
								},
							},
							"target": {
								Required: []string{
									"backend",
								},
								Properties: map[string]apiextensionsv1beta1.JSONSchemaProps{
									"path": {
										Description: "http path",
										Type:        "string",
										Pattern:     httpPathRegex,
										MinLength:   int64Ptr(1),
										Default: &apiextensionsv1beta1.JSON{
											Raw: []byte(""),
										},
									},
									"mock": {
										Description: "name of ConfigMap defining a valid mock response for this ApiProxy",
										Type:        "object",
										Required: []string{
											"configMapName",
										},
										Properties: map[string]apiextensionsv1beta1.JSONSchemaProps{
											"configMapName": {
												Description: "valid qualified Kubernetes name",
												Type:        "string",
												MinLength:   int64Ptr(1),
												MaxLength:   int64Ptr(63),
												Pattern:     k8sQualifiedNameRegex,
											},
										},
									},
									"backend": {
										Description: "defines the anatomy of an upstream server",
										OneOf: []apiextensionsv1beta1.JSONSchemaProps{
											{
												Description: "upstream server location outside a Kubernetes context",
												Type:        "object",
												Required: []string{
													"endpoint",
												},
												Properties: map[string]apiextensionsv1beta1.JSONSchemaProps{
													"endpoint": {
														Description: "endpoint object",
														Type:        "object",
                            Required: []string{
    													"scheme",
                              "host",
    												},
                            Properties: map[string]apiextensionsv1beta1.JSONSchemaProps{
                              "scheme": {
                                Description: "http url scheme",
                                Type:        "string",
                                Enum: []apiextensionsv1beta1.JSON{
    															{
    																Raw: []byte(`"http"`),
    															},
    															{
    																Raw: []byte(`"https"`),
    															},
    														},
                              },
                              "host": {
                                Description: "http url host",
                                Type:        "string",
                                MinLength:   int64Ptr(1),
                              },
                            },
													},
												},
											},
											{
												Description: "upstream server location inside a Kubernetes context",
												Type:        "object",
												Required: []string{
													"service",
												},
												Properties: map[string]apiextensionsv1beta1.JSONSchemaProps{
													"service": {
														Description: "dynamic or static Kubernetes service",
														Type:        "object",
														OneOf: []apiextensionsv1beta1.JSONSchemaProps{
															{
																Description: "statically defined Kubernetes service",
																Type:        "object",
																Required: []string{
																	"name",
																	"port",
																},
																Properties: map[string]apiextensionsv1beta1.JSONSchemaProps{
																	"name": {
																		Description: "valid qualified Kubernetes name",
																		Type:        "string",
																		MinLength:   int64Ptr(1),
																		MaxLength:   int64Ptr(63),
																		Pattern:     k8sQualifiedNameRegex,
																	},
																	"port": {
																		Description: "tcp port",
																		Type:        "integer",
																		Minimum:     float64Ptr(0),
																		Maximum:     float64Ptr(65535),
																	},
																},
															},
															{
																Description: "dynamically defined Kubernetes service",
																Type:        "object",
																Required: []string{
																	"labels",
																	"port",
																},
																Properties: map[string]apiextensionsv1beta1.JSONSchemaProps{
																	"port": {
																		Description: "tcp port",
																		Type:        "integer",
																		Minimum:     float64Ptr(0),
																		Maximum:     float64Ptr(65535),
																	},
																	"labels": {
																		Description: "list of Kubernetes service metadata labels to be matched against",
																		Type:        "array",
																		Items: &apiextensionsv1beta1.JSONSchemaPropsOrArray{
																			Schema: &apiextensionsv1beta1.JSONSchemaProps{
																				Description: "statically or dynamically defined label",
																				OneOf: []apiextensionsv1beta1.JSONSchemaProps{
																					{
																						Description: "statically defined metadata label",
																						Type:        "object",
																						Required: []string{
																							"name",
																							"value",
																						},
																						Properties: map[string]apiextensionsv1beta1.JSONSchemaProps{
																							"name": {
																								Description: "valid qualified Kubernetes name",
																								Type:        "string",
																								MinLength:   int64Ptr(1),
																								MaxLength:   int64Ptr(63),
																								Pattern:     k8sQualifiedNameRegex,
																							},
																							"value": {
																								Description: "service metadata label value",
																								Type:        "string",
																								Pattern:     k8sNameRegex,
																							},
																						},
																					},
																					{
																						Description: "dynamically defined metadata label based on http header value",
																						Type:        "object",
																						Required: []string{
																							"name",
																							"header",
																						},
																						Properties: map[string]apiextensionsv1beta1.JSONSchemaProps{
																							"name": {
																								Description: "valid qualified Kubernetes name",
																								Type:        "string",
																								MinLength:   int64Ptr(1),
																								MaxLength:   int64Ptr(63),
																								Pattern:     k8sQualifiedNameRegex,
																							},
																							"header": {
																								Description: "http header name",
																								Type:        "string",
																								Pattern:     httpHeaderNameRegex,
																								MinLength:   int64Ptr(1),
																							},
																						},
																					},
																				},
																			},
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
									"ssl": {
										Description: "kubernetes.io/tls secret type containing key, cert, and/or ca for ApiProxy tls configuration",
										Type:        "object",
										Required: []string{
											"secretName",
										},
										Properties: map[string]apiextensionsv1beta1.JSONSchemaProps{
											"secretName": {
												Description: "valid qualified Kubernetes name",
												Type:        "string",
												MinLength:   int64Ptr(1),
												MaxLength:   int64Ptr(63),
												Pattern:     k8sQualifiedNameRegex,
											},
										},
									},
								},
							},
							"plugins": {
								Description: "list of plugins to be executed on each request",
								Type:        "array",
								Items: &apiextensionsv1beta1.JSONSchemaPropsOrArray{
									Schema: &apiextensionsv1beta1.JSONSchemaProps{
										Description: "unique plugin item",
										Required: []string{
											"name",
										},
										Type: "object",
										Properties: map[string]apiextensionsv1beta1.JSONSchemaProps{
											"name": {
												Description: "plugin name",
												Type:        "string",
												MinLength:   int64Ptr(1),
											},
											"version": {
												Description: "plugin version",
												Type:        "string",
												Pattern:     semanticVersionRegex,
											},
											"config": {
												Description: "unstructured plugin configuration",
												Type:        "object",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	},
}

var MockTargetCRD = &apiextensionsv1beta1.CustomResourceDefinition{
	ObjectMeta: metav1.ObjectMeta{
		Name: fmt.Sprintf("mocktargets.%s", kanaliio.GroupName),
	},
	Spec: apiextensionsv1beta1.CustomResourceDefinitionSpec{
		Group:   kanaliio.GroupName,
		Version: version,
		Names: apiextensionsv1beta1.CustomResourceDefinitionNames{
			Plural:   "mocktargets",
			Singular: "mocktarget",
			ShortNames: []string{
				"mt",
			},
			Kind:     "MockTarget",
			ListKind: "MockTargetList",
		},
		Scope:      apiextensionsv1beta1.NamespaceScoped,
		Validation: &apiextensionsv1beta1.CustomResourceValidation{},
	},
}

func int64Ptr(f int64) *int64 {
	return &f
}

func stringPtr(f string) *string {
	return &f
}

func float64Ptr(f float64) *float64 {
	return &f
}
