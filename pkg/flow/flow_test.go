package flow

import (
	"context"
	"errors"
	"net/http"
	"testing"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/mocktracer"
	"github.com/stretchr/testify/assert"
)

type mockStep struct{}

func (s mockStep) Name() string {
	return "mock step"
}

func (s mockStep) Do(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	return nil
}

type mockErrorStep struct{}

func (s mockErrorStep) Name() string {
	return "mock error step"
}

func (s mockErrorStep) Do(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	return errors.New("forced error")
}

func TestNew(t *testing.T) {
	assert.Equal(t, New(), &Flow{})
}

func TestAdd(t *testing.T) {
	f := &Flow{}
	f.Add(mockStep{}, mockStep{})
	f.Add(mockStep{})
	assert.Equal(t, len(*f), 3)
	for _, step := range *f {
		assert.Equal(t, step.Name(), "mock step")
		assert.Nil(t, step.Do(context.Background(), nil, &http.Request{}))
	}
}

func TestPlay(t *testing.T) {
	f := &Flow{}
	f.Add(mockStep{})
	assert.Nil(t, f.Play(context.Background(), nil, &http.Request{}))
	f.Add(mockErrorStep{})
	sp := mocktracer.New().StartSpan("mock").(*mocktracer.MockSpan)
	assert.Error(t, f.Play(opentracing.ContextWithSpan(context.Background(), sp), nil, &http.Request{}))
	assert.Equal(t, len(sp.Tags()), 1)
	assert.Equal(t, len(sp.Logs()), 1)
}
