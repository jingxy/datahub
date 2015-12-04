package clog

import (
	"bytes"
	"os"
	"reflect"
	"testing"
)

func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	} else {
		t.Log("Test ok!")
	}
}

func TestLog(t *testing.T) {
	SetLogLevel(LOG_LEVEL_DEBUG)
	logstr := "hello world!!"
	Infof("%s", logstr)
	Debugf("%s", logstr)
	Errorf("%s", logstr)
	Warnf("%s", logstr)
	Tracef("%s", logstr)
	Printf("%s", logstr)
	Info(logstr)
	Debug(logstr)
	Error(logstr)
	Warn(logstr)
	Trace(logstr)
	Println(logstr)
	t.Log("TEST OK")

}

func TestLogLevel(t *testing.T) {
	SetLogLevel(LOG_LEVEL_TRACE)
	lvl := GetLogLevel()
	expect(t, lvl, LOG_LEVEL_TRACE)
}

func TestLogFile(t *testing.T) {
	SetLogFile("/tmp/asdsaaf")
	CloseLogFile()
	var null *os.File = nil
	expect(t, logfileFd, null)
}

func TestLogger(t *testing.T) {

	var buf bytes.Buffer

	s := "hello world!"
	SetOutput(&buf)
	Info(s)
	t.Log("buffer:", &buf)
	buf.Reset()
	Warn(s)
	t.Log("buffer:", &buf)
}

func TestLogLevelEnv(t *testing.T) {
	os.Setenv("DATAHUB_LOGLEVEL", "fatal")
	checkLogEnv()
	lvl := GetLogLevel()
	expect(t, lvl, LOG_LEVEL_FATAL)
}
