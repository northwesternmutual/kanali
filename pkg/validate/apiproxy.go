// Copyright (c) 2018 Northwestern Mutual.
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

package validate

import (
	"encoding/json"
	"errors"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/northwesternmutual/kanali/pkg/apis/kanali.io/v2"
	"github.com/northwesternmutual/kanali/pkg/log"
)

func (v *validation) IsValidApiProxy(data []byte) error {
	apiproxy := new(v2.ApiProxy)
	if err := json.Unmarshal(data, apiproxy); err != nil {
		return err
	}
	return v.isValidApiProxy(apiproxy)
}

// isValidApiProxy determines whether or not an ApiProxy can be applied.
//
// CASES
// 1. No ApiProxy with the same source.Path exists - VALID
// 2. An ApiProxy with the same source.Path exists
//    a. That ApiProxy is the same as the one being updated
//       i.  The source.VirtualHost is the same - VALID
//       ii. The source.VirtualHost is different
//           I.  The source.VirtualHost is equal to the empty string
//               1. There exists an ApiProxy with the same source.Path that has a non-empty source.VirtualHost - DENY
//               2. There does not exist an ApiProxy with the same source.Path that has a non-empty source.VirtualHost - VALID
//           II. The source.VirtualHost is non-empty
//               1. There exists another ApiProxy with the same source.Path and an empty source.VirtualHost - DENY
//               2. There exists an ApiProxy with the same source.Path and the same source.VirtualHost - DENY
//               3. There does not exist an ApiProxy with the same source.Path and source.VirtualHost- VALID
//    b. That ApiProxy is different than the one being updated
//       i.  The source.VirtualHost is equal to the empty string - DENY
//       ii. The source.VirtualHost is non-empty
//           I.   There exists another ApiProxy with the same source.Path and an empty source.VirtualHost - DENY
//           II.  There exists another ApiProxy with the same source.Path and source.VirtualHost - DENY
//           III. There does not exist another ApiProxy with the same source.Path and source.VirtualHost - VALID
func (v *validation) isValidApiProxy(apiproxy *v2.ApiProxy) error {
	logger := log.WithContext(v.ctx)

	// Start with a list of all ApiProxy resources currently in the cluster
	list, err := v.clientset.KanaliV2().ApiProxies(metav1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		return err
	}

	// BEGIN CASE 1: O(n), n => cardinality of ApiProxy resources
	var case1DirtyBit bool = true
	var case2DirtyBit *v2.ApiProxy
	for _, curr := range list.Items {
		if curr.Spec.Source.Path == apiproxy.Spec.Source.Path {
			case2DirtyBit = &curr
			case1DirtyBit = false
			break
		}
	}
	if case1DirtyBit {
		logger.Info(fmt.Sprintf("No ApiProxy with the same source.Path %s exists - VALID", apiproxy.Spec.Source.Path))
		return nil
	}
	// END CASE 1

	// BEGIN CASE 2.a
	if apiproxy.GetName() == case2DirtyBit.GetName() && apiproxy.GetNamespace() == case2DirtyBit.GetNamespace() {
		// BEGIN CASE 2.a.i: O(1)
		if apiproxy.Spec.Source.VirtualHost == case2DirtyBit.Spec.Source.VirtualHost {
			logger.Info("The existing ApiProxy is being updated. Soure path and virtual host is not changing")
			return nil
		}
		// END CASE 2.a.i

		// BEGIN CASE 2.a.ii
		if apiproxy.Spec.Source.VirtualHost != case2DirtyBit.Spec.Source.VirtualHost {
			// BEGIN CASE 2.a.ii.I
			if len(apiproxy.Spec.Source.VirtualHost) < 1 {
				for _, curr := range list.Items {
					if apiproxy.GetName() == curr.GetName() && apiproxy.GetNamespace() == curr.GetNamespace() {
						continue
					}
					if apiproxy.Spec.Source.Path == curr.Spec.Source.Path {
						// BEGIN CASE 2.a.ii.I.1
						if len(curr.Spec.Source.VirtualHost) > 0 {
							logger.Info("There exists an ApiProxy with the same source.Path that has a non-empty source.VirtualHost")
							return errors.New("There exists an ApiProxy with the same source.Path that has a non-empty source.VirtualHost")
						}
						// END CASE 2.a.ii.I.1
					}
				}
			}
			// END CASE 2.a.ii.I

			// BEGIN CASE 2.a.ii.II
			if len(apiproxy.Spec.Source.VirtualHost) > 0 {
				for _, curr := range list.Items {
					if apiproxy.GetName() == curr.GetName() && apiproxy.GetNamespace() == curr.GetNamespace() {
						continue
					}
					if apiproxy.Spec.Source.Path == curr.Spec.Source.Path {
						if curr.Spec.Source.VirtualHost == apiproxy.Spec.Source.VirtualHost {
							logger.Info("There exists another ApiProxy with the same source.Path and an empty source.VirtualHost")
							return errors.New("There exists another ApiProxy with the same source.Path and an empty source.VirtualHost")
						}
						return nil
					}
				}
			}
			// END CASE 2.a.ii.II
		}
		// END CASE 2.a.ii
	}
	// END CASE 2.a

	// BEGIN CASE 2.b
	if apiproxy.GetName() != case2DirtyBit.GetName() || apiproxy.GetNamespace() != case2DirtyBit.GetNamespace() {
		// BEGIN CASE 2.b.i
		if len(apiproxy.Spec.Source.VirtualHost) < 1 {
			logger.Info("The source.VirtualHost is equal to the empty string")
			return errors.New("The source.VirtualHost is equal to the empty string")
		}
		// END CASE 2.b.i

		// BEGIN CASE 2.b.ii
		if len(apiproxy.Spec.Source.VirtualHost) > 0 {
			for _, curr := range list.Items {
				if apiproxy.Spec.Source.Path == curr.Spec.Source.Path {
					if len(curr.Spec.Source.VirtualHost) < 1 {
						logger.Info("The source.VirtualHost is equal to the empty string")
						return errors.New("The source.VirtualHost is equal to the empty string")
					}

					if curr.Spec.Source.VirtualHost == apiproxy.Spec.Source.VirtualHost {
						logger.Info("There exists another ApiProxy with the same source.Path and source.VirtualHost")
						return errors.New("There exists another ApiProxy with the same source.Path and source.VirtualHost")
					}
				}
			}
		}
		// END CASE 2.b.ii
	}
	// END CASE 2.b

	return nil
}

func (v *validation) IsValidApiProxyList(data []byte) error {
	list := new(v2.ApiProxyList)
	if err := json.Unmarshal(data, list); err != nil {
		return err
	}
	return v.isValidApiProxyList(list)
}

func (v *validation) isValidApiProxyList(list *v2.ApiProxyList) error {
	for _, apiproxy := range list.Items {
		if err := v.isValidApiProxy(&apiproxy); err != nil {
			return err
		}
	}
	return nil
}
