package usecases

import (
	"fmt"
	"net"
	"os"
)

// StderrChatter - 標準エラー・チャッター。 標準エラー出力とロガーを一緒にしただけです。
type StderrChatter struct {
	logger Logger
}

// NewStderrChatter - チャッターを作成します。
func NewStderrChatter(logger Logger) *StderrChatter {
	chatter := new(StderrChatter)
	chatter.logger = logger
	return chatter
}

// Trace - 本番運用時にはソースコードにも残っていないような内容を書くのに使います。
func (chatter StderrChatter) Trace(text string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, text, args...) // 標準エラー出力
	chatter.logger.Trace(text, args...)   // ログ
}

// Debug - 本番運用時にもデバッグを取りたいような内容を書くのに使います。
func (chatter StderrChatter) Debug(text string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, text, args...)
	chatter.logger.Debug(text, args...)
}

// Info - 多めの情報を書くのに使います。
func (chatter StderrChatter) Info(text string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, text, args...)
	chatter.logger.Info(text, args...)
}

// Notice - 定期的に動作確認を取りたいような、節目、節目の重要なポイントの情報を書くのに使います。
func (chatter StderrChatter) Notice(text string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, text, args...)
	chatter.logger.Notice(text, args...)
}

// Warn - ハードディスクの残り容量が少ないなど、当面は無視できるが対応はしたいような情報を書くのに使います。
func (chatter StderrChatter) Warn(text string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, text, args...)
	chatter.logger.Warn(text, args...)
}

// Error - 動作不良の内容や、理由を書くのに使います。
func (chatter StderrChatter) Error(text string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, text, args...)
	chatter.logger.Error(text, args...)
}

// Fatal - 強制終了したことを伝えます。
func (chatter StderrChatter) Fatal(text string, args ...interface{}) string {
	message := fmt.Sprintf(text, args...)
	fmt.Fprintf(os.Stderr, message)
	chatter.logger.Fatal(message)
	return message
}

// Print - 必ず出力します。
func (chatter StderrChatter) Print(text string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, text, args...)
	chatter.logger.Print(text, args...) // ログ
}

// Send - メッセージを送信します。
func (chatter StderrChatter) Send(conn net.Conn, text string, args ...interface{}) {
	_, err := fmt.Fprintf(conn, text, args...) // 出力先指定
	if err != nil {
		panic(err)
	}

	fmt.Fprintf(os.Stderr, text, args...)
	chatter.logger.Print(text, args...)
}
