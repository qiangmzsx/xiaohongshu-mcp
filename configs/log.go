package configs

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	logDir    = ""
	logFile   *os.File
	logWriter io.Writer
)

func InitLogging() error {
	execPath, err := os.Executable()
	if err != nil {
		execPath, err = os.Getwd()
		if err != nil {
			logrus.Warnf("获取目录失败: %v, 使用默认日志配置", err)
			return nil
		}
	}

	logDir = filepath.Join(filepath.Dir(execPath), "logs")

	if err := os.MkdirAll(logDir, 0755); err != nil {
		logrus.Warnf("创建日志目录失败: %v, 使用默认日志配置", err)
		return nil
	}

	logFileName := time.Now().Format("xiaohongshu-mcp-2006-01-02-15.log")
	logFilePath := filepath.Join(logDir, logFileName)

	logFile, err = os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logrus.Warnf("打开日志文件失败: %v, 使用默认日志配置", err)
		return nil
	}

	logWriter = io.MultiWriter(os.Stdout, logFile)

	logrus.SetOutput(logWriter)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	logrus.Infof("日志初始化完成, 日志文件: %s", logFilePath)
	return nil
}

func GetLogDir() string {
	return logDir
}

func CloseLogging() {
	if logFile != nil {
		logFile.Close()
	}
}
