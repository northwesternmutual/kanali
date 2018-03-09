package run

import "context"

type Runnable interface {
	Run(context.Context) error
	Close(error) error
}

type Process struct {
	cancel context.CancelFunc
}

func MonitorContext(cancel context.CancelFunc) *Process {
	return &Process{
		cancel: cancel,
	}
}

func (p *Process) Run(ctx context.Context) error {
	<-ctx.Done()
	return nil
}

func (p *Process) Close(error) error {
	p.cancel()
	return nil
}
