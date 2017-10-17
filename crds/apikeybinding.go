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

var apiKeyBindingCRD = &apiextensionsv1beta1.CustomResourceDefinition{
	ObjectMeta: metav1.ObjectMeta{
		Name: fmt.Sprintf("apikeybindings.%s", KanaliGroupName),
	},
	Spec: apiextensionsv1beta1.CustomResourceDefinitionSpec{
		Group:   KanaliGroupName,
		Version: Version,
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
				AdditionalProperties: &apiextensionsv1beta1.JSONSchemaPropsOrBool{
					Allows: false,
				},
				Properties: map[string]apiextensionsv1beta1.JSONSchemaProps{
					"spec": {
						Description: "ApiKeyBinding resource body",
						Required: []string{
							"keys",
						},
						AdditionalProperties: &apiextensionsv1beta1.JSONSchemaPropsOrBool{
							Allows: false,
						},
						Properties: map[string]apiextensionsv1beta1.JSONSchemaProps{
							"keys": {
								Description: "list of ApiKey resources granted permissions",
								Type:        "array",
								UniqueItems: true,
								MinLength:   int64Ptr(1),
								Items: &apiextensionsv1beta1.JSONSchemaPropsOrArray{
									Schema: &apiextensionsv1beta1.JSONSchemaProps{
										Description: "ApiKey permissions",
										Type:        "object",
										Required: []string{
											"name",
										},
										AdditionalProperties: &apiextensionsv1beta1.JSONSchemaPropsOrBool{
											Allows: false,
										},
										Properties: map[string]apiextensionsv1beta1.JSONSchemaProps{
											"name": {
												Ref: stringPtr("#/definitions/name"),
											},
											"defaultRule": {
												Ref: stringPtr("#/definitions/rule"),
											},
											"subpaths": {
												Description: "list of subpaths defining fine grained permissions",
												Type:        "array",
												UniqueItems: true,
												MinLength:   int64Ptr(1),
												Items: &apiextensionsv1beta1.JSONSchemaPropsOrArray{
													Schema: &apiextensionsv1beta1.JSONSchemaProps{
														Description: "unique subpath",
														Required: []string{
															"path",
															"rule",
														},
														AdditionalProperties: &apiextensionsv1beta1.JSONSchemaPropsOrBool{
															Allows: false,
														},
														Properties: map[string]apiextensionsv1beta1.JSONSchemaProps{
															"path": {
																Ref: stringPtr("#/definitions/path"),
															},
															"rule": {
																Ref: stringPtr("#/definitions/rule"),
															},
														},
													},
												},
											},
											"quota": {
												Description: "number of requests authorized for this ApiKey",
												Ref:         stringPtr("#/definitions/wholeNumber"),
											},
											"rate": {
												Description: "number of requests an ApiKey can make over an interval",
												Type:        "object",
												Properties: map[string]apiextensionsv1beta1.JSONSchemaProps{
													"amount": {
														Ref: stringPtr("#/definitions/wholeNumber"),
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
				Definitions: map[string]apiextensionsv1beta1.JSONSchemaProps{
					"name": {
						Description: "valid qualified Kubernetes name",
						Type:        "string",
						MinLength:   int64Ptr(1),
						MaxLength:   int64Ptr(63),
						Pattern:     k8sQualifiedNameRegex,
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
					"wholeNumber": {
						Type:       "integer",
						Minimum:    float64Ptr(1),
						MultipleOf: float64Ptr(1),
					},
					"rule": {
						Description: "defines http method that an ApiKey has access to",
						Type:        "object",
						OneOf: []apiextensionsv1beta1.JSONSchemaProps{
							{
								Description: "global access",
								Type:        "object",
								AdditionalProperties: &apiextensionsv1beta1.JSONSchemaPropsOrBool{
									Allows: false,
								},
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
								AdditionalProperties: &apiextensionsv1beta1.JSONSchemaPropsOrBool{
									Allows: false,
								},
								Required: []string{
									"granular",
								},
								Properties: map[string]apiextensionsv1beta1.JSONSchemaProps{
									"granular": {
										Description: "fine grained http verb access",
										Type:        "object",
										AdditionalProperties: &apiextensionsv1beta1.JSONSchemaPropsOrBool{
											Allows: false,
										},
										Required: []string{
											"verbs",
										},
										Properties: map[string]apiextensionsv1beta1.JSONSchemaProps{
											"verbs": {
												Description: "list of http methods that ApiKey has access to",
												Type:        "array",
												UniqueItems: true,
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
}
