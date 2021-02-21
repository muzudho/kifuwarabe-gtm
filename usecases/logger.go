package usecases

import (
	"bufio"
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

	traceFile  *os.File
	debugFile  *os.File
	infoFile   *os.File
	noticeFile *os.File
	warnFile   *os.File
	errorFile  *os.File
	fatalFile  *os.File
	printFile  *os.File

	traceWriter  *bufio.Writer
	debugWriter  *bufio.Writer
	infoWriter   *bufio.Writer
	noticeWriter *bufio.Writer
	warnWriter   *bufio.Writer
	errorWriter  *bufio.Writer
	fatalWriter  *bufio.Writer
	printWriter  *bufio.Writer
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

// OpenAllLogs - 全てのログ・ファイルを開けます
func (logger *Logger) OpenAllLogs() error {
	file, err := logger.openLogFile(logger.tracePath)
	if err != nil {
		return err
	}
	logger.traceFile = file
	logger.traceWriter = bufio.NewWriter(file)

	file, err = logger.openLogFile(logger.debugPath)
	if err != nil {
		return err
	}
	logger.debugFile = file
	logger.debugWriter = bufio.NewWriter(file)

	file, err = logger.openLogFile(logger.infoPath)
	if err != nil {
		return err
	}
	logger.infoFile = file
	logger.infoWriter = bufio.NewWriter(file)

	file, err = logger.openLogFile(logger.noticePath)
	if err != nil {
		return err
	}
	logger.noticeFile = file
	logger.noticeWriter = bufio.NewWriter(file)

	file, err = logger.openLogFile(logger.warnPath)
	if err != nil {
		return err
	}
	logger.warnFile = file
	logger.warnWriter = bufio.NewWriter(file)

	file, err = logger.openLogFile(logger.errorPath)
	if err != nil {
		return err
	}
	logger.errorFile = file
	logger.errorWriter = bufio.NewWriter(file)

	file, err = logger.openLogFile(logger.fatalPath)
	if err != nil {
		return err
	}
	logger.fatalFile = file
	logger.fatalWriter = bufio.NewWriter(file)

	file, err = logger.openLogFile(logger.printPath)
	if err != nil {
		return err
	}
	logger.printFile = file
	logger.printWriter = bufio.NewWriter(file)

	return nil
}

// OpenAllLogs - 全てのログ・ファイルを開けます
func (logger *Logger) openLogFile(filePath string) (*os.File, error) {
	// 追加書込み。
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		// ログのファイル・オープン失敗
		fmt.Fprintf(os.Stderr, "filePath=[%s]\n", filePath)
		fmt.Fprintf(os.Stderr, "err=[%s]\n", err)
		return nil, err
	}
	return file, nil
}

// FlushAllLogs - バッファーに溜まっている分をファイルに書き出します。定期的に行ってください
func (logger *Logger) FlushAllLogs() {
	logger.traceWriter.Flush()
	logger.debugWriter.Flush()
	logger.infoWriter.Flush()
	logger.noticeWriter.Flush()
	logger.warnWriter.Flush()
	logger.errorWriter.Flush()
	logger.fatalWriter.Flush()
	logger.printWriter.Flush()
}

// CloseAllLogs - 全てのログ・ファイルを閉じます
func (logger *Logger) CloseAllLogs() {
	logger.FlushAllLogs()

	defer logger.traceFile.Close()
	defer logger.debugFile.Close()
	defer logger.infoFile.Close()
	defer logger.noticeFile.Close()
	defer logger.warnFile.Close()
	defer logger.errorFile.Close()
	defer logger.fatalFile.Close()
	defer logger.printFile.Close()
}

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
func (logger *Logger) write(file *os.File, text string, args ...interface{}) {
	// tはtime.Time型
	t := time.Now()

	s := fmt.Sprintf(text, args...)
	s = fmt.Sprintf("[%s] %s", t.Format(timeStampLayout), s)
	fmt.Fprint(file, s)
}

// Trace - ログファイルに書き込みます。
func (logger Logger) Trace(text string, args ...interface{}) {
	logger.traceWriter.WriteString(fmt.Sprintf(text, args...))
}

// Debug - ログファイルに書き込みます。
func (logger Logger) Debug(text string, args ...interface{}) {
	logger.debugWriter.WriteString(fmt.Sprintf(text, args...))
}

// Info - ログファイルに書き込みます。
func (logger Logger) Info(text string, args ...interface{}) {
	logger.infoWriter.WriteString(fmt.Sprintf(text, args...))
}

// Notice - ログファイルに書き込みます。
func (logger Logger) Notice(text string, args ...interface{}) {
	logger.noticeWriter.WriteString(fmt.Sprintf(text, args...))
}

// Warn - ログファイルに書き込みます。
func (logger Logger) Warn(text string, args ...interface{}) {
	logger.warnWriter.WriteString(fmt.Sprintf(text, args...))
}

// Error - ログファイルに書き込みます。
func (logger Logger) Error(text string, args ...interface{}) {
	logger.errorWriter.WriteString(fmt.Sprintf(text, args...))
}

// Fatal - ログファイルに書き込みます。
func (logger Logger) Fatal(text string, args ...interface{}) string {
	message := fmt.Sprintf(text, args...)
	logger.fatalWriter.WriteString(message)
	return message
}

// Print - ログファイルに書き込みます。 Chatter から呼び出してください。
func (logger Logger) Print(text string, args ...interface{}) {
	logger.printWriter.WriteString(fmt.Sprintf(text, args...))
}
