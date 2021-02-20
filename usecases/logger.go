package usecases

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Logger - ロガー。
type Logger struct {
	tracePath  string
	debugPath  string
	infoPath   string
	noticePath string
	warnPath   string
	errorPath  string
	fatalPath  string
	printPath  string
}

// NewLogger - ロガーを作成します。
func NewLogger(
	tracePath string,
	debugPath string,
	infoPath string,
	noticePath string,
	warnPath string,
	errorPath string,
	fatalPath string,
	printPath string) *Logger {

	logger := new(Logger)
	logger.tracePath = tracePath
	logger.debugPath = debugPath
	logger.infoPath = infoPath
	logger.noticePath = noticePath
	logger.warnPath = warnPath
	logger.errorPath = errorPath
	logger.fatalPath = fatalPath
	logger.printPath = printPath

	return logger
}

// Go言語では、 yyyy とかではなく、定められた数をそこに置くのらしい☆（＾～＾）
const timeStampLayout = "2006-01-02 15:04:05"

// RemoveAllOldLogs - 既存のログファイルを削除します
// 誤動作防止のため、 basename の末尾が '.log' か、または basename に '.log.' が含まれるものだけ削除できるものとします。
func (logger *Logger) RemoveAllOldLogs() {
	logger.removeLog(logger.tracePath)
	logger.removeLog(logger.debugPath)
	logger.removeLog(logger.infoPath)
	logger.removeLog(logger.noticePath)
	logger.removeLog(logger.warnPath)
	logger.removeLog(logger.errorPath)
	logger.removeLog(logger.fatalPath)
	logger.removeLog(logger.printPath)
}

// 誤動作防止のため、 basename の末尾が '.log' か、または basename に '.log.' が含まれるものだけ削除できるものとします。
func (logger *Logger) removeLog(path string) {
	basename := filepath.Base(path)
	if strings.HasSuffix(basename, ".log") || strings.Contains(basename, ".log.") {
		os.Remove(path)
	}
}

// write - ログファイルに書き込みます
func write(filePath string, text string, args ...interface{}) {
	// TODO ファイルの開閉回数を減らせないものか。
	// 追加書込み。
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	// tはtime.Time型
	t := time.Now()

	s := fmt.Sprintf(text, args...)
	s = fmt.Sprintf("[%s] %s", t.Format(timeStampLayout), s)
	fmt.Fprint(file, s)
	defer file.Close()
}

// Trace - ログファイルに書き込みます。
func (logger Logger) Trace(text string, args ...interface{}) {
	write(logger.tracePath, text, args...)
}

// Debug - ログファイルに書き込みます。
func (logger Logger) Debug(text string, args ...interface{}) {
	write(logger.debugPath, text, args...)
}

// Info - ログファイルに書き込みます。
func (logger Logger) Info(text string, args ...interface{}) {
	write(logger.infoPath, text, args...)
}

// Notice - ログファイルに書き込みます。
func (logger Logger) Notice(text string, args ...interface{}) {
	write(logger.noticePath, text, args...)
}

// Warn - ログファイルに書き込みます。
func (logger Logger) Warn(text string, args ...interface{}) {
	write(logger.warnPath, text, args...)
}

// Error - ログファイルに書き込みます。
func (logger Logger) Error(text string, args ...interface{}) {
	write(logger.errorPath, text, args...)
}

// Fatal - ログファイルに書き込みます。
func (logger Logger) Fatal(text string, args ...interface{}) string {
	message := fmt.Sprintf(text, args...)
	write(logger.fatalPath, message)
	return message
}

// Print - ログファイルに書き込みます。 Chatter から呼び出してください。
func (logger Logger) Print(text string, args ...interface{}) {
	write(logger.printPath, text, args...)
}
