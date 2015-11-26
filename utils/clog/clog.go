package clog

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
)

var (
	KNRM = "\x1B[0m"
	KBLD = "\x1B[1m"
	KITY = "\x1B[3m"
	KUND = "\x1B[4m"
	KRED = "\x1B[31m"
	KGRN = "\x1B[32m"
	KYEL = "\x1B[33m"
	KBLU = "\x1B[34m"
	KMAG = "\x1B[35m"
	KCYN = "\x1B[36m"
	KWHT = "\x1B[37m"
)

const (
	LOG_LEVEL_NONE = iota
	LOG_LEVEL_FATAL
	LOG_LEVEL_ERROR
	LOG_LEVEL_WARN
	LOG_LEVEL_INFO
	LOG_LEVEL_TRACE
	LOG_LEVEL_DEBUG
)

var logEnv = map[string]int{
	"none":  LOG_LEVEL_NONE,
	"fatal": LOG_LEVEL_FATAL,
	"error": LOG_LEVEL_ERROR,
	"warn":  LOG_LEVEL_WARN,
	"info":  LOG_LEVEL_INFO,
	"trace": LOG_LEVEL_TRACE,
	"debug": LOG_LEVEL_DEBUG,
}

var logLevel = LOG_LEVEL_INFO
var logfileFd *os.File

func trace() string {
	pc := make([]uintptr, 5) // at least 1 entry needed
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[1]).Name()

	fName := strings.Split(f, "/")[strings.Count(f, "/")]

	//return fmt.Sprintf("%s() ", f)
	return fName + "() "
}

func SetLogLevel(level int) {
	logLevel = level
}

func GetLogLevel() (level int) {
	return logLevel
}

func SetLogFile(logfile string) {
	f, err := os.OpenFile(logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	} else {
		log.SetOutput(f)
		logfileFd = f
	}
}

func CloseLogFile() {
	if logfileFd != nil {
		logfileFd.Close()
		logfileFd = nil
	}
}

func SetOutput(w io.Writer) {
	log.SetOutput(w)
}

func Error(a ...interface{}) {
	checkLogEnv()
	if logLevel >= LOG_LEVEL_ERROR {
		log.Print(KRED+"[ERROR] "+KNRM+trace(), fmt.Sprintln(a...))
	}
}

func Errorf(format string, a ...interface{}) {
	checkLogEnv()
	if logLevel >= LOG_LEVEL_ERROR {
		log.Print(KRED+"[ERROR] "+KNRM+trace(), fmt.Sprintf(format, a...))
	}
}

func Fatal(a ...interface{}) {
	checkLogEnv()
	if logLevel >= LOG_LEVEL_FATAL {
		log.Print(KRED+KBLD+"[FATAL] "+KNRM+trace(), fmt.Sprintln(a...))
		os.Exit(1)
	}
}

func Fatalf(format string, a ...interface{}) {
	checkLogEnv()
	if logLevel >= LOG_LEVEL_FATAL {
		log.Print(KRED+KBLD+"[FATAL] "+KNRM+trace(), fmt.Sprintf(format, a...))
		os.Exit(1)
	}
}

func Info(a ...interface{}) {
	checkLogEnv()
	if logLevel >= LOG_LEVEL_INFO {
		log.Print(KGRN+"[INFO] "+KNRM+trace(), fmt.Sprintln(a...))
	}
}

func Infof(format string, a ...interface{}) {
	checkLogEnv()
	if logLevel >= LOG_LEVEL_INFO {
		log.Print(KGRN+"[INFO] "+KNRM+trace(), fmt.Sprintf(format, a...))
	}
}

func Trace(a ...interface{}) {
	checkLogEnv()
	if logLevel >= LOG_LEVEL_TRACE {
		log.Print(KMAG+"[TRACE] "+KNRM+trace(), fmt.Sprintln(a...))
	}
}

func Tracef(format string, a ...interface{}) {
	checkLogEnv()
	if logLevel >= LOG_LEVEL_TRACE {
		log.Print(KMAG+"[TRACE] "+KNRM+trace(), fmt.Sprintf(format, a...))
	}
}
func Debug(a ...interface{}) {
	checkLogEnv()
	if logLevel >= LOG_LEVEL_DEBUG {
		log.Print(KBLU+"[DEBUG] "+KNRM+trace(), fmt.Sprintln(a...))
	}
}

func Debugf(format string, a ...interface{}) {
	checkLogEnv()
	if logLevel >= LOG_LEVEL_DEBUG {
		log.Print(KBLU+"[DEBUG] "+KNRM+trace(), fmt.Sprintf(format, a...))
	}
}

func Warn(a ...interface{}) {
	checkLogEnv()
	if logLevel >= LOG_LEVEL_WARN {
		log.Print(KYEL+"[WARNING] "+KNRM+trace(), fmt.Sprintln(a...))
	}
}

func Warnf(format string, a ...interface{}) {
	checkLogEnv()
	if logLevel >= LOG_LEVEL_WARN {
		log.Print(KYEL+"[WARNING] "+KNRM+trace(), fmt.Sprintf(format, a...))
	}
}

func Println(a ...interface{}) {
	Info(a...)
}

func Printf(format string, a ...interface{}) {
	Infof(format, a...)
}

func checkLogEnv() {
	lvl := os.Getenv("DATAHUB_LOGLEVEL")
	if len(lvl) > 0 {
		lvl = strings.ToLower(lvl)
		if level, ok := logEnv[lvl]; ok {
			if logLevel != level {
				fmt.Printf("set log level to %v[%v]\n", lvl, level)
				logLevel = level
			}
		}
	}
}

func init() {
	checkLogEnv()
}

/*
func main() {
	fmt.Println(defaultLogLevel)
	test()
	Info("%s", "hello world!")
	Debug("%s", "hello world!")
	Error("%s", "hello world!")
	Warn("%s", "hello world!")
	Fatal("%s", "hello world!")
}

func init() {
	//log.SetFlags(log.Lshortfile | log.LstdFlags)
	fmt.Println(LOG_LEVEL_DEBUG, LOG_LEVEL_ERROR, LOG_LEVEL_FATAL, LOG_LEVEL_INFO, LOG_LEVEL_TRACE, LOG_LEVEL_WARN)
}
*/
