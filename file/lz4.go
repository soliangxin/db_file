package file

import (
	"github.com/pierrec/lz4"
	"io"
)

const (
	DefaultLz4CompressionLevel = 0
)

type Lz4Write struct {
	w *lz4.Writer
}

func (l *Lz4Write) NewWriter(w io.Writer) (err error) {
	l.w = lz4.NewWriter(w)
	l.w.Header = lz4.Header{
		CompressionLevel: DefaultLz4CompressionLevel,
	}
	return err
}

func (l *Lz4Write) Write(p []byte) (int, error) {
	return l.w.Write(p)
}

func (l *Lz4Write) WriteString(s string) (int, error) {
	return l.w.Write([]byte(s))
}

func (l *Lz4Write) Reset(w io.Writer) {
	l.w.Reset(w)
}

func (l *Lz4Write) Flush() error {
	return l.w.Flush()
}

func (l *Lz4Write) Close() error {
	return l.w.Close()
}
