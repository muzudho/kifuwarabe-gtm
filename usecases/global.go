package usecases

const (
	// Author - 囲碁思考エンジンの作者名だぜ☆（＾～＾）
	Author = "Satoshi Takahashi"
)

// GlobalVariables - グローバル変数。
type GlobalVariables struct {
	// Log - ロガー。
	Log Logger
	// Chat - チャッター。 標準出力とロガーを一緒にしただけです。
	Chat Chatter
	// StderrChat - チャッター。 標準エラー出力とロガーを一緒にしただけです。
	StderrChat StderrChatter
}

// G - グローバル変数。思い切った名前。
var G GlobalVariables
