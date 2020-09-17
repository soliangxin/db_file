package file

import (
	"bytes"
	gzip "github.com/klauspost/pgzip"
	"io/ioutil"
	"testing"
)

func TestGzipWrite(t *testing.T) {
	testStr := "test hello"
	var buf bytes.Buffer

	// 写入数据
	gw := GzipWrite{}
	err := gw.NewWriter(&buf)
	if err != nil {
		t.Error(err)
	}
	_, err = gw.Write([]byte(testStr))
	if err != nil {
		t.Error(err)
	}
	err = gw.Close()
	if err != nil {
		t.Error(err)
	}

	// 读取数据
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
}
