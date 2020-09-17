package db_file

import (
	"bytes"
	"github.com/golang/snappy"
	gzip "github.com/klauspost/pgzip"
	"github.com/pierrec/lz4"
	"github.com/valyala/gozstd"
	"io/ioutil"
	"testing"
)

func TestOpenCmpFileWriteGzip(t *testing.T) {
	// 写入数据
	testStr := "test hello"
	var buf bytes.Buffer

	// 写入数据
	gw, err := NewWrite(&buf, "gzip")
	if err != nil {
		t.Error("new CompressWrite gzip error, error message: ", err)
	}
	_, err = gw.Write([]byte(testStr))
	if err != nil {
		t.Error(err)
	}
	err = gw.Close()
	if err != nil {
		t.Error(err)
	}

	gr, err := gzip.NewReader(&buf)
	if err != nil {
		t.Error(err)
	}
	data, err := ioutil.ReadAll(gr)
	if err != nil {
		t.Error(err)
	}
	if string(data[:]) != testStr {
		t.Errorf("gzip un error, expect: %s, current: %s", testStr, string(data[:]))
	}
	err = gr.Close()
	if err != nil {
		t.Error(err)
	}
}

func TestOpenCmpFileWriteLz4(t *testing.T) {
	// 写入数据
	testStr := "test hello"
	var buf bytes.Buffer

	// 写入数据
	gw, err := NewWrite(&buf, "lz4")
	if err != nil {
		t.Error("new CompressWrite lz4 error, error message: ", err)
	}
	_, err = gw.Write([]byte(testStr))
	if err != nil {
		t.Error(err)
	}
	err = gw.Close()
	if err != nil {
		t.Error(err)
	}

	gr := lz4.NewReader(&buf)
	if err != nil {
		t.Error(err)
	}
	data, err := ioutil.ReadAll(gr)
	if err != nil {
		t.Error(err)
	}
	if string(data[:]) != testStr {
		t.Errorf("lz4 un error, expect: %s, current: %s", testStr, string(data[:]))
	}
}

func TestOpenCmpFileWriteSnappy(t *testing.T) {
	// 写入数据
	testStr := "test hello"
	var buf bytes.Buffer

	// 写入数据
	gw, err := NewWrite(&buf, "snappy")
	if err != nil {
		t.Error("new CompressWrite snappy error, error message: ", err)
	}
	_, err = gw.Write([]byte(testStr))
	if err != nil {
		t.Error(err)
	}
	err = gw.Close()
	if err != nil {
		t.Error(err)
	}

	gr := snappy.NewReader(&buf)
	if err != nil {
		t.Error(err)
	}
	data, err := ioutil.ReadAll(gr)
	if err != nil {
		t.Error(err)
	}
	if string(data[:]) != testStr {
		t.Errorf("snappy un error, expect: %s, current: %s", testStr, string(data[:]))
	}
}

func TestOpenCmpFileWriteZstd(t *testing.T) {
	// 写入数据
	testStr := "test hello"
	var buf bytes.Buffer

	// 写入数据
	gw, err := NewWrite(&buf, "zstd")
	if err != nil {
		t.Error("new CompressWrite zstd error, error message: ", err)
	}
	_, err = gw.Write([]byte(testStr))
	if err != nil {
		t.Error(err)
	}
	err = gw.Close()
	if err != nil {
		t.Error(err)
	}

	gr := gozstd.NewReader(&buf)
	if err != nil {
		t.Error(err)
	}
	data, err := ioutil.ReadAll(gr)
	if err != nil {
		t.Error(err)
	}
	if string(data[:]) != testStr {
		t.Errorf("zstd un error, expect: %s, current: %s", testStr, string(data[:]))
	}
}
