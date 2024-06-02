package log

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

const (
	allLog  = "all"
	errLog  = "err"
	warnLog = "warn"
	infoLog = "info"
)

type FileDateHook struct {
	file     *os.File
	errFile  *os.File
	warnFile *os.File
	infoFile *os.File
	logPath  string
	fileDate string
}

func (hook FileDateHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
func (hook FileDateHook) Fire(entry *logrus.Entry) error {
	timer := entry.Time.Format("2006-01-02_15-04")
	line, _ := entry.String()
	if hook.fileDate == timer {
		switch entry.Level {
		case logrus.ErrorLevel:
			hook.errFile.Write([]byte(line))
		case logrus.WarnLevel:
			hook.warnFile.Write([]byte(line))
		case logrus.InfoLevel:
			hook.infoFile.Write([]byte(line))
		}
		hook.file.Write([]byte(line))
		return nil
	}
	hook.file.Close()
	os.MkdirAll(fmt.Sprintf("%s/%s", hook.logPath, timer), os.ModePerm)
	hook.file, _ = os.OpenFile(fmt.Sprintf("%s/%s/%s.log", hook.logPath, timer, allLog), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
	hook.errFile, _ = os.OpenFile(fmt.Sprintf("%s/%s/%s.log", hook.logPath, timer, errLog), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
	hook.warnFile, _ = os.OpenFile(fmt.Sprintf("%s/%s/%s.log", hook.logPath, timer, warnLog), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
	hook.infoFile, _ = os.OpenFile(fmt.Sprintf("%s/%s/%s.log", hook.logPath, timer, infoLog), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)

	switch entry.Level {
	case logrus.ErrorLevel:
		hook.errFile.Write([]byte(line))
	case logrus.WarnLevel:
		hook.warnFile.Write([]byte(line))
	case logrus.InfoLevel:
		hook.infoFile.Write([]byte(line))
	}
	hook.fileDate = timer
	hook.file.Write([]byte(line))
	return nil
}

type MyFormatter struct {
}

func (f MyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer == nil {
		b = &bytes.Buffer{}
	} else {
		b = entry.Buffer
	}
	// 设置格式
	fmt.Fprintf(b, "%s\n", entry.Message)

	return b.Bytes(), nil
}

func InitFile(logPath string) {
	logrus.SetFormatter(&MyFormatter{})

	fileDate := time.Now().Format("2006-01-02_15-04")
	//创建目录
	err := os.MkdirAll(fmt.Sprintf("%s/%s", logPath, fileDate), os.ModePerm)
	if err != nil {
		logrus.Error(err)
		return
	}
	allFile, err := os.OpenFile(fmt.Sprintf("%s/%s/%s.log", logPath, fileDate, allLog), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
	errFile, err := os.OpenFile(fmt.Sprintf("%s/%s/%s.log", logPath, fileDate, errLog), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
	warnFile, err := os.OpenFile(fmt.Sprintf("%s/%s/%s.log", logPath, fileDate, warnLog), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
	infoFile, err := os.OpenFile(fmt.Sprintf("%s/%s/%s.log", logPath, fileDate, infoLog), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)

	fileHook := FileDateHook{allFile, errFile, warnFile, infoFile, logPath, fileDate}

	logrus.AddHook(&fileHook)
}
