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

// This file was automatically generated by informer-gen

package v2

import (
	kanali_io_v2 "github.com/northwesternmutual/kanali/pkg/apis/kanali.io/v2"
	versioned "github.com/northwesternmutual/kanali/pkg/client/clientset/versioned"
	internalinterfaces "github.com/northwesternmutual/kanali/pkg/client/informers/externalversions/internalinterfaces"
	v2 "github.com/northwesternmutual/kanali/pkg/client/listers/kanali/v2"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
	time "time"
)

// MockTargetInformer provides access to a shared informer and lister for
// MockTargets.
type MockTargetInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v2.MockTargetLister
}

type mockTargetInformer struct {
	factory internalinterfaces.SharedInformerFactory
}

// NewMockTargetInformer constructs a new informer for MockTarget type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewMockTargetInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				return client.KanaliV2().MockTargets(namespace).List(options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				return client.KanaliV2().MockTargets(namespace).Watch(options)
			},
		},
		&kanali_io_v2.MockTarget{},
		resyncPeriod,
		indexers,
	)
}

func defaultMockTargetInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewMockTargetInformer(client, v1.NamespaceAll, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc})
}

func (f *mockTargetInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&kanali_io_v2.MockTarget{}, defaultMockTargetInformer)
}

func (f *mockTargetInformer) Lister() v2.MockTargetLister {
	return v2.NewMockTargetLister(f.Informer().GetIndexer())
}