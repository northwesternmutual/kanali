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

  metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
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
        Required: []string{
          "spec",
        },
        AdditionalProperties: &apiextensionsv1beta1.JSONSchemaPropsOrBool{
          Allows: false,
        },
        Properties: map[string]apiextensionsv1beta1.JSONSchemaProps{
          "keys": apiextensionsv1beta1.JSONSchemaProps{
            Type: "array",
            UniqueItems: true,
            MinLength: int64Ptr(1),
            Items: &apiextensionsv1beta1.JSONSchemaPropsOrArray{
              Schema: &apiextensionsv1beta1.JSONSchemaProps{
                Type: "object",
                Required: []string{
                  "data",
                  "status",
                },
                Properties: map[string]apiextensionsv1beta1.JSONSchemaProps{
                  "data": apiextensionsv1beta1.JSONSchemaProps{
                    Type: "string",
                  },
                  "status": apiextensionsv1beta1.JSONSchemaProps{
                    Type: "string",
                    Enum: []apiextensionsv1beta1.JSON{
      								{
      									Raw: []byte(`"active"`),
      								},
      								{
      									Raw: []byte(`"inactive"`),
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
