package chain

import (
	"net/http"
)

// Chain represents a set of middlewares
type Chain []func(next http.Handler) http.Handler

// New instantiates a new chain.
func New() *Chain {
	return &Chain{}
}

// Add inserts a new middleware into the chain.
func (c *Chain) Add(m ...func(next http.Handler) http.Handler) *Chain {
	for _, v := range m {
		if v == nil {
			continue
		}
		*c = append(*c, v)
	}
	return c
}

// Link will link together off of the middlewares
// in the chain. The resulting http.Handler will
// execute the middleware in the order they were
// inserted into the chain with the passed middleware
// executed last.
func (c *Chain) Link(final http.HandlerFunc) http.Handler {
	if c == nil || final == nil {
		return nil
	}
	switch len(*c) {
	case 0:
		return nil
	case 1:
		return (*c)[0](final)
	}
	return (*c).link(final)
}

func (c Chain) link(inner http.Handler) http.Handler {
	if len(c) > 1 {
		inner = c[1:].link(inner)
	}
	return c[0](inner)
}
