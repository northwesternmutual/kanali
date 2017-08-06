package flow

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/northwesternmutual/kanali/controller"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/assert"
)

type mockStep struct{}

func (s mockStep) GetName() string {
	return "mock step"
}

func (s mockStep) Do(ctx context.Context, ctlr *controller.Controller, w http.ResponseWriter, r *http.Request, resp *http.Response, trace opentracing.Span) error {
	return nil
}

type mockErrorStep struct{}

func (s mockErrorStep) GetName() string {
	return "mock error step"
}

func (s mockErrorStep) Do(ctx context.Context, ctlr *controller.Controller, w http.ResponseWriter, r *http.Request, resp *http.Response, trace opentracing.Span) error {
	return errors.New("forced error")
}

func TestAdd(t *testing.T) {
	f := &Flow{}
	f.Add(mockStep{}, mockStep{})
	f.Add(mockStep{})
	assert.Equal(t, len(*f), 3)
	for _, step := range *f {
		assert.Equal(t, step.GetName(), "mock step")
		assert.Nil(t, step.Do(context.Background(), nil, nil, nil, nil, opentracing.StartSpan("test span")))
	}
}

func TestPlay(t *testing.T) {
	f := &Flow{}
	f.Add(mockStep{})
	assert.Nil(t, f.Play(context.Background(), nil, nil, nil, nil, opentracing.StartSpan("test span")))
	f.Add(mockErrorStep{})
	assert.Error(t, f.Play(context.Background(), nil, nil, nil, nil, opentracing.StartSpan("test span")))
}
