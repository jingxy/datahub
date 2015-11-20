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

func Error(format string, a ...interface{}) {
	log.Println(KRED+"[ERROR] "+KNRM+trace(), fmt.Sprintf(format, a...))
}
func Fatal(format string, a ...interface{}) {
	log.Println(KRED+KBLD+"[FATAL] "+KNRM+trace(), fmt.Sprintf(format, a...))
	os.Exit(1)
}

func Info(format string, a ...interface{}) {
	log.Println(KGRN+"[INFO] "+KNRM+trace(), fmt.Sprintf(format, a...))
}

func Debug(format string, a ...interface{}) {
	log.Println(KBLU+"[DEBUG] "+KNRM+trace(), fmt.Sprintf(format, a...))
}

func Warn(format string, a ...interface{}) {
	log.Println(KYEL+"[WARNING] "+KNRM+trace(), fmt.Sprintf(format, a...))
}

func Printf(format string, a ...interface{}) {
	log.Println(KGRN+"[INFO] "+KNRM+trace(), fmt.Sprintf(format, a...))
}

func Println(a ...interface{}) {
	log.Printf(KGRN+"[INFO] "+KNRM+trace()+"%s", fmt.Sprintln(a...))
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
