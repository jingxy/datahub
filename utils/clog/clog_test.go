package clog

import (
	"bytes"
	"fmt"
	"testing"
)

func TestLog(t *testing.T) {
	Infof("%s", "hello world!")
	Debugf("%s", "hello world!")
	Errorf("%s", "hello world!")
	Warnf("%s", "hello world!")
	t.Log("...")

}

/*

func TestLogger(t *testing.T) {

	var buf bytes.Buffer

	cloger := New(&buf)
	s := "hello world! warn!!!"
	//cloger.SetOutput(&buf)
	cloger.Log(s)
	fmt.Print("clogger:", &buf)
}
*/
