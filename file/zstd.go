package file

import (
	"github.com/valyala/gozstd"
	"io"
)

const (
	DefaultZstdCompressionLevel = gozstd.DefaultCompressionLevel
)

type ZstdWrite struct {
	w *gozstd.Writer
}

func (z *ZstdWrite) NewWriter(w io.Writer) (err error) {
	z.w = gozstd.NewWriterLevel(w, DefaultZstdCompressionLevel)
	return nil
}

func (z *ZstdWrite) Write(p []byte) (int, error) {
	return z.w.Write(p)
}

func (z *ZstdWrite) WriteString(s string) (int, error) {
	return z.w.Write([]byte(s))
}

func (z *ZstdWrite) Reset(w io.Writer) {
	n := &gozstd.CDict{}
	z.w.Reset(w, n, DefaultZstdCompressionLevel)
}

func (z *ZstdWrite) Flush() error {
	return z.w.Flush()
}

func (z *ZstdWrite) Close() error {
	return z.w.Close()
}
