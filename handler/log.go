package handler

import (
	"bufio"
	"fmt"
	"log"
	"matching/config"
	"os"
	"path/filepath"
	"runtime"
	_ "strings"
	"sync"
	"time"
)

type LEVEL byte

// 定义几种日志级别
const (
	DEBUG LEVEL = iota
	INFO
	WARN
	ERROR
)
var LevelArr = []LEVEL{DEBUG, INFO, WARN, ERROR}

const DateFormat = "2006-01-02"

type FileLogger struct {
	fileDir        string         // 日志文件保存的目录
	fileName       string         // 日志文件名（无需包含日期和扩展名）
	prefix         string         // 日志消息的前缀
	//logLevel       LEVEL          // 日志等级
	lfDebug        *os.File       // 日志文件
	lfInfo         *os.File       // 日志文件
	lfWarn         *os.File       // 日志文件
	lfError        *os.File       // 日志文件
	bufDebug       *bufio.Writer    // 日志缓冲
	bufInfo        *bufio.Writer    // 日志缓冲
	bufWarn        *bufio.Writer    // 日志缓冲
	bufError       *bufio.Writer    // 日志缓冲
	date           *time.Time     // 日志当前日期
	mu             *sync.RWMutex  // 读写锁，在进行日志分割和日志写入时需要锁住
	logChan        chan LogContent// 日志消息通道，以实现异步写日志
	stopTickerChan chan bool      // 停止定时器的通道
}

type LogContent struct {
	level LEVEL
	content string
}

var fileLogger *FileLogger

// 初始化函数
func Init() error {

	// 判断日志开关
	if !config.LogSwitch {
		return nil
	}

	CloseLogger()

	f := &FileLogger{
		fileDir:        "",
		fileName:       "file",
		prefix:         "",
		mu:             new(sync.RWMutex),
		logChan:        make(chan LogContent, 5000),
		stopTickerChan: make(chan bool, 1),
	}

	nowpath, _ := os.Getwd()
	f.fileDir = nowpath + "/logs/"

	t, _ := time.Parse(DateFormat, time.Now().Format(DateFormat))
	f.date = &t

	// 判断文件和文件夹是否存在
	f.isExistOrCreateFileDir()
	err := f.isExistOrCreateFile()
	if err != nil {
		return err
	}

	// 开几个协程处理日志
	go f.logWriter()
	go f.fileMonitor()

	fileLogger = f

	return nil
}

// 关闭日志
func CloseLogger() {
	if fileLogger != nil {
		fileLogger.stopTickerChan <- true
		close(fileLogger.stopTickerChan)
		close(fileLogger.logChan)
		fileLogger.bufDebug = nil
		fileLogger.bufInfo = nil
		fileLogger.bufWarn = nil
		fileLogger.bufError = nil
		fileLogger.lfDebug.Close()
		fileLogger.lfInfo.Close()
		fileLogger.lfWarn.Close()
		fileLogger.lfError.Close()
	}
}

// 判断文件夹是否存在，不存在则创建
func (f *FileLogger) isExistOrCreateFileDir() {
	var fileDir string
	for v := range LevelArr {
		switch v {
		case 0:
			fileDir = filepath.Join(f.fileDir, "debug/")
		case 1:
			fileDir = filepath.Join(f.fileDir, "info/")
		case 2:
			fileDir = filepath.Join(f.fileDir, "warn/")
		case 3:
			fileDir = filepath.Join(f.fileDir, "error/")
		}

		if _, err := os.Stat(fileDir); os.IsNotExist(err) {
			// 必须分成两步
			// 先创建文件夹
			os.Mkdir(fileDir, 0755)
			// 再修改权限
			os.Chmod(fileDir, 0755)
		}
	}
}

// 判断文件是否存在，不存在就创建
func (f *FileLogger) isExistOrCreateFile() error {

	var level string
	//var prefix string
	for v := range LevelArr {
		switch v {
		case 0:
			level = "debug"
			//prefix = "[Debug]"
			// 对每个文件进行创建读取
			filePath := filepath.Join(f.fileDir, level, f.fileName)
			file := filePath + "-" + f.date.Format(DateFormat) + ".log"
			// 这里文件不需要判断了，直接用open自动创建模式
			lfDebug, err := os.OpenFile(file, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
			if err != nil {
				log.Println("文件创建失败：", file)
				return err
			}
			f.lfDebug = lfDebug
			f.bufDebug = bufio.NewWriter(lfDebug)
		case 1:
			level = "info"
			//prefix = "[Info]"
			// 对每个文件进行创建读取
			filePath := filepath.Join(f.fileDir, level, f.fileName)
			file := filePath + "-" + f.date.Format(DateFormat) + ".log"
			// 这里文件不需要判断了，直接用open自动创建模式
			lfInfo, err := os.OpenFile(file, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
			if err != nil {
				log.Println("文件创建失败：", file)
				return err
			}
			f.lfInfo = lfInfo
			f.bufInfo = bufio.NewWriter(lfInfo)
		case 2:
			level = "warn"
			//prefix = "[Warn]"
			// 对每个文件进行创建读取
			filePath := filepath.Join(f.fileDir, level, f.fileName)
			file := filePath + "-" + f.date.Format(DateFormat) + ".log"
			// 这里文件不需要判断了，直接用open自动创建模式
			lfWarn, err := os.OpenFile(file, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
			if err != nil {
				log.Println("文件创建失败：", file)
				return err
			}
			f.lfWarn = lfWarn
			f.bufWarn = bufio.NewWriter(lfWarn)
		case 3:
			level = "error"
			//prefix = "[Error]"
			// 对每个文件进行创建读取
			filePath := filepath.Join(f.fileDir, level, f.fileName)
			file := filePath + "-" + f.date.Format(DateFormat) + ".log"
			// 这里文件不需要判断了，直接用open自动创建模式
			lfError, err := os.OpenFile(file, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
			if err != nil {
				log.Println("文件创建失败：", file)
				return err
			}
			f.lfError = lfError
			f.bufError = bufio.NewWriter(lfError)
		}
	}
	return nil
}

// 非堵塞写入日志
func (f *FileLogger) logWriter() {
	defer func() { recover() }()
	for {
		content, ok := <-f.logChan
		if !ok {
			return
		}

		f.mu.RLock()
		// 根据日志类型判断写入哪个文件
		switch content.level {
		case DEBUG:
			//f.lfDebug.Write([]byte(content.content))
			f.bufDebug.WriteString(content.content)
		case INFO:
			//f.lgInfo.Output(2, content.content)
		case WARN:
			//f.lgWarn.Output(2, content.content)
		case ERROR:
			//f.lgError.Output(2, content.content)
		}

		f.mu.RUnlock()
	}
}

// 文件监听，每30秒检查文件是否需要分隔
func (f *FileLogger) fileMonitor() {
	defer func() { recover() }()
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			log.Println("检查文件分隔")
			if f.isMustSplit() {
				if err := f.split(); err != nil {
					log.Printf("Log split error: %v\n", err)
				}
			}
		case <-f.stopTickerChan:
			return
		}
	}
}

// 默认必须要分隔
func (f *FileLogger) isMustSplit() bool {
	return true
}

// 分隔文件
func (f *FileLogger) split() error {
	f.mu.Lock() // 先加锁
	defer f.mu.Unlock()

	// 先关闭文件连接
	if f.lfDebug != nil {
		f.lfDebug.Close()
	}
	if f.lfInfo != nil {
		f.lfDebug.Close()
	}
	if f.lfWarn != nil {
		f.lfDebug.Close()
	}
	if f.lfError != nil {
		f.lfDebug.Close()
	}

	// 重新设置时间
	t, _ := time.Parse(DateFormat, time.Now().Format(DateFormat))
	f.date = &t

	// 创建文件并重新绑定文件句柄
	err := f.isExistOrCreateFile()
	if err != nil {
		return err
	}

	return nil
}

// 对外提供的接口
func LogDebug(format string, v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	var logContent LogContent
	logContent.level = 0
	logContent.content = fmt.Sprintf("[%v:%v]", filepath.Base(file), line) + fmt.Sprintf("[Debug]"+format, v...)
	fileLogger.logChan <- logContent
}

func LogInfo(format string, v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	var logContent LogContent
	logContent.level = 1
	logContent.content = fmt.Sprintf("[%v:%v]", filepath.Base(file), line) + fmt.Sprintf("[Info]"+format, v...)
	fileLogger.logChan <- logContent
}