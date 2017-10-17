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

package crds

import (
	"fmt"

	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	httpPathRegex = `^\/.*`
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

var apiProxyCRD = &apiextensionsv1beta1.CustomResourceDefinition{
	ObjectMeta: metav1.ObjectMeta{
		Name: fmt.Sprintf("apiproxies.%s", KanaliGroupName),
	},
	Spec: apiextensionsv1beta1.CustomResourceDefinitionSpec{
		Group:   KanaliGroupName,
		Version: Version,
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
				AdditionalProperties: &apiextensionsv1beta1.JSONSchemaPropsOrBool{
					Allows: false,
				},
				Properties: map[string]apiextensionsv1beta1.JSONSchemaProps{
					"spec": {
						Description: "ApiProxy resource body",
						Required: []string{
							"source",
							"target",
						},
						AdditionalProperties: &apiextensionsv1beta1.JSONSchemaPropsOrBool{
							Allows: false,
						},
						Properties: map[string]apiextensionsv1beta1.JSONSchemaProps{
							"source": {
								Description: "unique incoming http path and virtual host combination",
								Required: []string{
									"path",
								},
								AdditionalProperties: &apiextensionsv1beta1.JSONSchemaPropsOrBool{
									Allows: false,
								},
								Properties: map[string]apiextensionsv1beta1.JSONSchemaProps{
									"path": {
										Ref: stringPtr("#/definitions/path"),
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
								AdditionalProperties: &apiextensionsv1beta1.JSONSchemaPropsOrBool{
									Allows: false,
								},
								Properties: map[string]apiextensionsv1beta1.JSONSchemaProps{
									"path": {
										Ref: stringPtr("#/definitions/path"),
									},
									"mock": {
										Description: "name of ConfigMap defining a valid mock response for this ApiProxy",
										Type:        "object",
										Required: []string{
											"configMapName",
										},
										AdditionalProperties: &apiextensionsv1beta1.JSONSchemaPropsOrBool{
											Allows: false,
										},
										Properties: map[string]apiextensionsv1beta1.JSONSchemaProps{
											"configMapName": {
												Ref: stringPtr("#/definitions/name"),
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
												AdditionalProperties: &apiextensionsv1beta1.JSONSchemaPropsOrBool{
													Allows: false,
												},
												Properties: map[string]apiextensionsv1beta1.JSONSchemaProps{
													"endpoint": {
														Description: "valid http url",
														Type:        "string",
														Pattern:     httpURLRegex,
													},
												},
											},
											{
												Description: "upstream server location inside a Kubernetes context",
												Type:        "object",
												Required: []string{
													"service",
												},
												AdditionalProperties: &apiextensionsv1beta1.JSONSchemaPropsOrBool{
													Allows: false,
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
																AdditionalProperties: &apiextensionsv1beta1.JSONSchemaPropsOrBool{
																	Allows: false,
																},
																Properties: map[string]apiextensionsv1beta1.JSONSchemaProps{
																	"name": {
																		Ref: stringPtr("#/definitions/name"),
																	},
																	"port": {
																		Ref: stringPtr("#/definitions/port"),
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
																AdditionalProperties: &apiextensionsv1beta1.JSONSchemaPropsOrBool{
																	Allows: false,
																},
																Properties: map[string]apiextensionsv1beta1.JSONSchemaProps{
																	"port": {
																		Ref: stringPtr("#/definitions/port"),
																	},
																	"labels": {
																		Description: "list of Kubernetes service metadata labels to be matched against",
																		Type:        "array",
																		UniqueItems: true,
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
																						AdditionalProperties: &apiextensionsv1beta1.JSONSchemaPropsOrBool{
																							Allows: false,
																						},
																						Properties: map[string]apiextensionsv1beta1.JSONSchemaProps{
																							"name": {
																								Ref: stringPtr("#/definitions/name"),
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
																						AdditionalProperties: &apiextensionsv1beta1.JSONSchemaPropsOrBool{
																							Allows: false,
																						},
																						Properties: map[string]apiextensionsv1beta1.JSONSchemaProps{
																							"name": {
																								Ref: stringPtr("#/definitions/name"),
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
										AdditionalProperties: &apiextensionsv1beta1.JSONSchemaPropsOrBool{
											Allows: false,
										},
										Properties: map[string]apiextensionsv1beta1.JSONSchemaProps{
											"secretName": {
												Ref: stringPtr("#/definitions/name"),
											},
										},
									},
								},
							},
							"plugins": {
								Description: "list of plugins to be executed on each request",
								Type:        "array",
								UniqueItems: true,
								Items: &apiextensionsv1beta1.JSONSchemaPropsOrArray{
									Schema: &apiextensionsv1beta1.JSONSchemaProps{
										Description: "unique plugin item",
										Required: []string{
											"name",
										},
										AdditionalProperties: &apiextensionsv1beta1.JSONSchemaPropsOrBool{
											Allows: false,
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
				Definitions: map[string]apiextensionsv1beta1.JSONSchemaProps{
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
					"path": {
						Description: "http path",
						Type:        "string",
						Pattern:     httpPathRegex,
						MinLength:   int64Ptr(1),
						Default: &apiextensionsv1beta1.JSON{
							Raw: []byte("/"),
						},
					},
				},
			},
		},
	},
}
