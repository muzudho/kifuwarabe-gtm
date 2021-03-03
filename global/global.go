package global

import (
	l "github.com/muzudho/go-logger"
)

const (
	// Author - 囲碁思考エンジンの作者名だぜ☆（＾～＾）
	Author = "Satoshi Takahashi"
)

// Variables - グローバル変数。
type Variables struct {
	// Log - ロガー。
	Log l.Logger
	// Chat - チャッター。 標準出力とロガーを一緒にしただけです。
	Chat l.Chatter
	// StderrChat - チャッター。 標準エラー出力とロガーを一緒にしただけです。
	StderrChat l.StderrChatter
}

// G - グローバル変数。思い切った名前。
var G Variables
