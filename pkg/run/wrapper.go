package run

import (
	"context"

	"k8s.io/client-go/tools/cache"
)

type informerWrapper struct {
	informer cache.SharedInformer
}

func InformerWrapper(c cache.SharedInformer) Runnable {
	return &informerWrapper{
		informer: c,
	}
}

func (w *informerWrapper) Run(ctx context.Context) error {
	w.informer.Run(ctx.Done())
	return nil
}

func (w *informerWrapper) Close(error) error {
	return nil
}
