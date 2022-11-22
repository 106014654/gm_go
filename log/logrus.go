package log

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

//日志自定义格式
type LogFormatter struct{}

//格式详情
func (s *LogFormatter) Format(entry *log.Entry) ([]byte, error) {
	timestamp := time.Now().Local().Format(time.RFC3339)
	var file string
	var len int
	if entry.Caller != nil {
		file = filepath.Base(entry.Caller.File)
		len = entry.Caller.Line
	}
	//fmt.Println(entry.Data)
	msg := fmt.Sprintf("%s [%s:%d][%s] %sn", timestamp, file, len, strings.ToUpper(entry.Level.String()), entry.Message)
	return []byte(msg), nil
}

func getGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

type logFileWriter struct {
	file     *os.File
	logPath  string
	fileName string //判断日期切换目录

}

func (p *logFileWriter) Write(data []byte) (n int, err error) {
	if p == nil {
		return 0, errors.New("logFileWriter is nil")
	}
	if p.file == nil {
		return 0, errors.New("file not opened")
	}

	FilterLogDate(p.logPath, p.fileName)

	n, e := p.file.Write(data)
	return n, e
}

func FilterLogDate(logPath string, fileName string) {
	fileDate := time.Now().Format("2006-01-11")

	filename := fmt.Sprintf("%s/%s.log", logPath, fileName+fileDate)
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE|os.O_SYNC, 0600)
	if err != nil {
		log.Error(err)
		return
	}
	fileWriter := logFileWriter{file, logPath, fileName}
	log.SetOutput(&fileWriter)
}

//初始化日志
func InitLog(logPath string, fileName string) {
	FilterLogDate(logPath, fileName)

	log.SetReportCaller(true)
	log.SetFormatter(new(LogFormatter))
}
