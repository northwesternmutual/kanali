package run

import (
	"context"

	"github.com/oklog/run"
)

type Group struct {
	group run.Group
}

func (g *Group) Add(ctx context.Context, r Runnable) {
	g.group.Add(func() error {
		return r.Run(ctx)
	}, func(err error) {
		r.Close(err)
	})
}

func (g *Group) Run() error {
	return g.group.Run()
}
