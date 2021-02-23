package global

import u "github.com/muzudho/kifuwarabe-gtp/usecases"

const (
	// Author - 囲碁思考エンジンの作者名だぜ☆（＾～＾）
	Author = "Satoshi Takahashi"
)

// GlobalVariables - グローバル変数。
type GlobalVariables struct {
	// Log - ロガー。
	Log u.Logger
	// Chat - チャッター。 標準出力とロガーを一緒にしただけです。
	Chat u.Chatter
	// StderrChat - チャッター。 標準エラー出力とロガーを一緒にしただけです。
	StderrChat u.StderrChatter
}

// G - グローバル変数。思い切った名前。
var G GlobalVariables
