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
