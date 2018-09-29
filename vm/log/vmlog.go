
package vmlog

import (
	"fmt"
	"log"
	"os"
	"time"
)

const (
	//0
	PanicLevel int = iota
	//1
	FatalLevel
	//2
	ErrorLevel
	WarnLevel
	InfoLevel
	//5
	DebugLevel
)


type VmLog struct {
	level    int
	logTime  int64
	fileName string
	//file fd
	fileFd   *os.File

}
//
//var logger   *log
//
var logFile VmLog
// config log
func Config(logFolder string, level int) {
	logFile.fileName = logFolder
	logFile.level = level
	logFile.createLogFile()
	log.SetOutput(logFile.fileFd)
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)
	return
}

func DebugPrint(format string, args ...interface{}) {
	if logFile.level >= DebugLevel {
		log.SetPrefix("[DEBUG] ")
		log.Output(2, fmt.Sprintf(format, args...))
		//out to std
		fmt.Sprintf(format, args...)
	}
}

func InfoPrint(format string, args ...interface{}) {
	if logFile.level >= InfoLevel {
		log.SetPrefix("[INFO] ")
		log.Output(2, fmt.Sprintf(format, args...))
	}
}

func WarnPrint(format string, args ...interface{}) {
	if logFile.level >= WarnLevel {
		log.SetPrefix("[WARN] ")
		log.Output(2, fmt.Sprintf(format, args...))
	}
}

func ErrorPrint(format string, args ...interface{}) {
	if logFile.level >= ErrorLevel {
		log.SetPrefix("[ERROR] ")
		log.Output(2, fmt.Sprintf(format, args...))
	}
}

func FatalPrint(format string, args ...interface{}) {
	if logFile.level >= FatalLevel {
		log.SetPrefix("[FATAL] ")
		log.Output(2, fmt.Sprintf(format, args...))
	}
}

//log write to file
//func (lf VmLog) Write(buf []byte) (n int, err error) {
//	if lf.fileName == "" {
//		fmt.Printf("consol: %s", buf)
//		return len(buf), nil
//	}
//
//	if logFile.logTime+3600 < time.Now().Unix() {
//		logFile.createLogFile()
//		logFile.logTime = time.Now().Unix()
//	}
//
//	if logFile.fileFd == nil {
//		return len(buf), nil
//	}
//
//	return logFile.fileFd.Write(buf)
//}
//create log file
func (lf *VmLog) createLogFile() {

	now := time.Now()
	filename := fmt.Sprintf("%s_%04d%02d%02d_%02d%02d%02d.log", lf.fileName, now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())

	if fd, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModeExclusive); nil == err {
		lf.fileFd.Sync()
		lf.fileFd.Close()
		lf.fileFd = fd
		return
	}


	lf.fileFd = nil
	return
}
