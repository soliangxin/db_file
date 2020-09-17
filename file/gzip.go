package file

import (
	gzip "github.com/klauspost/pgzip"
	"io"
)

const (
	DefaultGzipCompressionLevel = gzip.DefaultCompression
)

type GzipWrite struct {
	w *gzip.Writer
}

func (g *GzipWrite) NewWriter(w io.Writer) (err error) {
	g.w, err = gzip.NewWriterLevel(w, DefaultGzipCompressionLevel)
	return err
}

func (g *GzipWrite) Write(p []byte) (int, error) {
	return g.w.Write(p)
}

func (g *GzipWrite) WriteString(s string) (int, error) {
	return g.w.Write([]byte(s))
}

func (g *GzipWrite) Reset(w io.Writer) {
	g.w.Reset(w)
}

func (g *GzipWrite) Flush() error {
	return g.w.Flush()
}

func (g *GzipWrite) Close() error {
	return g.w.Close()
}
