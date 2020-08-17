package dry

import (
	"compress/flate"
	"compress/gzip"
	"io"
	"sync"
)

var (
	Deflate DeflatePool
	Gzip    GzipPool
)

// DeflatePool manages a pool of flate.Writer
// flate.NewWriter allocates a lot of memory, so if flate.Writer
// are needed frequently, it's more efficient to use a pool of them.
// The pool uses sync.Pool internally.
type DeflatePool struct {
	pool sync.Pool
}

// GetWriter returns flate.Writer from the pool, or creates a new one
// with flate.BestCompression if the pool is empty.
func (pool *DeflatePool) GetWriter(dst io.Writer) (writer *flate.Writer) {
	if w := pool.pool.Get(); w != nil {
		writer = w.(*flate.Writer)
		writer.Reset(dst)
	} else {
		writer, _ = flate.NewWriter(dst, flate.BestCompression)
	}
	return writer
}

// ReturnWriter returns a flate.Writer to the pool that can
// late be reused via GetWriter.
// Don't close the writer, Flush will be called before returning
// it to the pool.
func (pool *DeflatePool) ReturnWriter(writer *flate.Writer) {
	writer.Close()
	pool.pool.Put(writer)
}

// GzipPool manages a pool of gzip.Writer.
// The pool uses sync.Pool internally.
type GzipPool struct {
	pool sync.Pool
}

// GetWriter returns gzip.Writer from the pool, or creates a new one
// with gzip.BestCompression if the pool is empty.
func (pool *GzipPool) GetWriter(dst io.Writer) (writer *gzip.Writer) {
	if w := pool.pool.Get(); w != nil {
		writer = w.(*gzip.Writer)
		writer.Reset(dst)
	} else {
		writer, _ = gzip.NewWriterLevel(dst, gzip.BestCompression)
	}
	return writer
}

// ReturnWriter returns a gzip.Writer to the pool that can
// late be reused via GetWriter.
// Don't close the writer, Flush will be called before returning
// it to the pool.
func (pool *GzipPool) ReturnWriter(writer *gzip.Writer) {
	writer.Close()
	pool.pool.Put(writer)
}
