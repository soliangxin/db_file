package file

import (
	"bytes"
	"github.com/golang/snappy"
	"io/ioutil"
	"testing"
)

func TestSnappyWrite(t *testing.T) {
	testStr := "test hello"
	var buf bytes.Buffer

	// 写入数据
	gw := SnappyWrite{}
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
	gr := snappy.NewReader(&buf)
	data, err := ioutil.ReadAll(gr)
	if err != nil {
		t.Error(err)
	}
	if string(data[:]) != testStr {
		t.Errorf("snappy un error, expect: %s, current: %s", testStr, string(data[:]))
	}
}
