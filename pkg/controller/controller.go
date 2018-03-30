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

package controller

import (
	"crypto/rsa"

	"github.com/northwesternmutual/kanali/pkg/client/informers/externalversions"
	"github.com/northwesternmutual/kanali/pkg/controller/apikey"
	"github.com/northwesternmutual/kanali/pkg/controller/apikeybinding"
	"github.com/northwesternmutual/kanali/pkg/controller/apiproxy"
	"github.com/northwesternmutual/kanali/pkg/controller/mocktarget"
)

func InitEventHandlers(f externalversions.SharedInformerFactory, decryptionKey *rsa.PrivateKey) {
	kanaliV2Informer := f.Kanali().V2()

	kanaliV2Informer.ApiKeys().Informer().AddEventHandler(apikey.NewController(decryptionKey))
	kanaliV2Informer.ApiKeyBindings().Informer().AddEventHandler(apikeybinding.NewController())
	kanaliV2Informer.ApiProxies().Informer().AddEventHandler(apiproxy.NewController())
	kanaliV2Informer.MockTargets().Informer().AddEventHandler(mocktarget.NewController())
}
