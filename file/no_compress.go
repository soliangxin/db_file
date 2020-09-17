package file

import (
	"bufio"
	"io"
)

type NoCompression struct {
	w *bufio.Writer
}

func (l *NoCompression) NewWriter(w io.Writer) (err error) {
	l.w = bufio.NewWriter(w)
	return err
}

func (l *NoCompression) Write(p []byte) (int, error) {
	return l.w.Write(p)
}

func (l *NoCompression) WriteString(s string) (int, error) {
	return l.w.Write([]byte(s))
}

func (l *NoCompression) Reset(w io.Writer) {
	l.w.Reset(w)
}

func (l *NoCompression) Flush() error {
	return l.w.Flush()
}

func (l *NoCompression) Close() error {
	return l.w.Flush()
}
