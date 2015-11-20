package clog

import (
	"fmt"
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

func trace() string {
	pc := make([]uintptr, 5) // at least 1 entry needed
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[1]).Name()

	fName := strings.Split(f, "/")[strings.Count(f, "/")]

	//return fmt.Sprintf("%s() ", f)
	return fName + "() "
}

func Error(a ...interface{}) {
	log.Print(KRED+"[ERROR] "+KNRM+trace(), fmt.Sprintln(a...))
}

func Errorf(format string, a ...interface{}) {
	log.Print(KRED+"[ERROR] "+KNRM+trace(), fmt.Sprintf(format, a...))
}

func Fatal(a ...interface{}) {
	log.Print(KRED+KBLD+"[FATAL] "+KNRM+trace(), fmt.Sprintln(a...))
	os.Exit(1)
}

func Fatalf(format string, a ...interface{}) {
	log.Print(KRED+KBLD+"[FATAL] "+KNRM+trace(), fmt.Sprintf(format, a...))
	os.Exit(1)
}

func Info(a ...interface{}) {
	log.Print(KGRN+"[INFO] "+KNRM+trace(), fmt.Sprintln(a...))
}

func Infof(format string, a ...interface{}) {
	log.Print(KGRN+"[INFO] "+KNRM+trace(), fmt.Sprintf(format, a...))
}

func Trace(a ...interface{}) {
	log.Print(KWHT+"[TRACE] "+KNRM+trace(), fmt.Sprintln(a...))
}

func Tracef(format string, a ...interface{}) {
	log.Print(KWHT+"[TRACE] "+KNRM+trace(), fmt.Sprintf(format, a...))
}
func Debug(a ...interface{}) {
	log.Print(KBLU+"[DEBUG] "+KNRM+trace(), fmt.Sprintln(a...))
}

func Debugf(format string, a ...interface{}) {
	log.Print(KBLU+"[DEBUG] "+KNRM+trace(), fmt.Sprintf(format, a...))
}

func Warn(a ...interface{}) {
	log.Print(KYEL+"[WARNING] "+KNRM+trace(), fmt.Sprintln(a...))
}


func Warnf(format string, a ...interface{}) {
	log.Print(KYEL+"[WARNING] "+KNRM+trace(), fmt.Sprintf(format, a...))
}

func Println(a ...interface{}) {
	log.Print(KGRN+"[INFO] "+KNRM+trace(), fmt.Sprintln(a...))
}

func Printf(format string, a ...interface{}) {
	log.Print(KGRN+"[INFO] "+KNRM+trace(), fmt.Sprintf(format, a...))
}


/*
func test() {
	Warn("test...")
}

func main() {
	test()
	Info("%s", "hello world!")
	Debug("%s", "hello world!")
	Error("%s", "hello world!")
	Fatal("%s", "hello world!")
	Warn("%s", "hello world!")
	Pf("%s", "sdsd")
	Pfln("hello")
}
*/
func init() {
	//log.SetFlags(log.Lshortfile | log.LstdFlags)
}
