package file

import (
	"github.com/golang/snappy"
	"io"
)

type SnappyWrite struct {
	w *snappy.Writer
}

func (s *SnappyWrite) NewWriter(w io.Writer) (err error) {
	s.w = snappy.NewBufferedWriter(w)
	return err
}

func (s *SnappyWrite) Write(p []byte) (int, error) {
	return s.w.Write(p)
}

func (s *SnappyWrite) WriteString(st string) (int, error) {
	return s.w.Write([]byte(st))
}

func (s *SnappyWrite) Reset(w io.Writer) {
	s.w.Reset(w)
}

func (s *SnappyWrite) Flush() error {
	return s.w.Flush()
}

func (s *SnappyWrite) Close() error {
	return s.w.Close()
}
