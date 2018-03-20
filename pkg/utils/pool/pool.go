package pool

import (
	"bytes"
	"sync"
)

// Read more about sync.Pool here:
// https://golang.org/pkg/sync/#Pool
var buffers = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

// GetBuffer fetches a buffer from the pool
func GetBuffer() *bytes.Buffer {
	return buffers.Get().(*bytes.Buffer)
}

// PutBuffer returns a buffer to the pool
func PutBuffer(buf *bytes.Buffer) {
	buf.Reset()
	buffers.Put(buf)
}
