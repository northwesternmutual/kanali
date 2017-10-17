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
	rfc3339Regex       = `^([0-9]+)-(0[1-9]|1[012])-(0[1-9]|[12][0-9]|3[01])[Tt]([01][0-9]|2[0-3]):([0-5][0-9]):([0-5][0-9]|60)(\.[0-9]+)?(([Zz])|([\+|\-]([01][0-9]|2[0-3]):[0-5][0-9]))$`
	encryptedDataRegex = `[0-9a-zA-Z]+`
)

var apiKeyCRD = &apiextensionsv1beta1.CustomResourceDefinition{
	ObjectMeta: metav1.ObjectMeta{
		Name: fmt.Sprintf("apikeys.%s", KanaliGroupName),
	},
	Spec: apiextensionsv1beta1.CustomResourceDefinitionSpec{
		Group:   KanaliGroupName,
		Version: Version,
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
				AdditionalProperties: &apiextensionsv1beta1.JSONSchemaPropsOrBool{
					Allows: false,
				},
				Properties: map[string]apiextensionsv1beta1.JSONSchemaProps{
					"spec": {
						Description: "ApiKey resource body",
						Required: []string{
							"revisions",
						},
						AdditionalProperties: &apiextensionsv1beta1.JSONSchemaPropsOrBool{
							Allows: false,
						},
						Properties: map[string]apiextensionsv1beta1.JSONSchemaProps{
							"revisions": {
								Description: "represents the list of active and inactive API key revisions",
								Type:        "array",
								UniqueItems: true,
								MinLength:   int64Ptr(1),
								Items: &apiextensionsv1beta1.JSONSchemaPropsOrArray{
									Schema: &apiextensionsv1beta1.JSONSchemaProps{
										Description: "a represetation of an API key revision",
										Type:        "object",
										Required: []string{
											"data",
											"status",
										},
										AdditionalProperties: &apiextensionsv1beta1.JSONSchemaPropsOrBool{
											Allows: false,
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
														Raw: []byte(`"active"`),
													},
													{
														Raw: []byte(`"inactive"`),
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
