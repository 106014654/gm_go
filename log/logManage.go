package log

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type logformatter struct{}

// 格式详情
func (s *logformatter) Format(entry *log.Entry) ([]byte, error) {
	timestamp := time.Now().Local().Format("2006-01-02 15:04:05")
	var file string
	var len int
	if entry.Caller != nil {
		file = filepath.Base(entry.Caller.File)
		len = entry.Caller.Line
	}
	//fmt.println(entry.data)
	msg := fmt.Sprintf("%s [%s:%d][%s] %s\n", timestamp, file, len, strings.ToUpper(entry.Level.String()), entry.Message)
	return []byte(msg), nil
}

type logfilewriter struct {
	file     *os.File
	logpath  string
	filedate string //判断日期切换目录
	appname  string
}

func (p *logfilewriter) Write(data []byte) (n int, err error) {
	if p == nil {
		return 0, errors.New("logfilewriter is nil")
	}
	if p.file == nil {
		return 0, errors.New("file not opened")
	}

	//判断是否需要切换日期
	filedate := time.Now().Format("20060102")
	if p.filedate != filedate {
		p.file.Close()

		filename := fmt.Sprintf("%s/%s-%s.log", p.logpath, p.appname, filedate)

		p.file, err = os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE|os.O_SYNC, 0600)
		if err != nil {
			return 0, err
		}

	}

	n, e := p.file.Write(data)
	return n, e

}

// 初始化日志
func Initlog(logpath string, appname string) {
	filedate := time.Now().Format("20060102")

	filename := fmt.Sprintf("%s/%s-%s.log", logpath, appname, filedate)

	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE|os.O_SYNC, 0600)

	if err != nil {
		log.Error(err)
		return
	}

	filewriter := logfilewriter{file, logpath, filedate, appname}
	log.SetOutput(&filewriter)

	log.SetReportCaller(true)
	log.SetFormatter(new(logformatter))
}
