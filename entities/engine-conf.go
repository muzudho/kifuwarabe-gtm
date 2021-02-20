package entities

import (
	"strings"
	"unicode"

	"github.com/muzudho/kifuwarabe-gtp/entities/stone"
)

// EngineConf - Tomlファイル。
type EngineConf struct {
	Profile Profile
	Engine  Engine
}

// Profile - [Profile] 区画。
type Profile struct {
	// Name - 対局者名（アカウント名）
	// Only A-Z a-z 0-9
	// Names may be at most 10 characters long
	Name string
	Pass string
}

// Engine - [Engine] テーブル。
type Engine struct {
	Komi      float32
	BoardSize int8
	MaxMoves  int16
	BoardData string
}

func removeAllWhiteSpace(str string) string {
	var b strings.Builder
	b.Grow(len(str))
	for _, ch := range str {
		if !unicode.IsSpace(ch) {
			b.WriteRune(ch)
		}
	}
	return b.String()
}

// GetBoardArray - 盤上の石の色の配列。
func (config EngineConf) GetBoardArray() []int {
	nodes := removeAllWhiteSpace(config.Engine.BoardData)
	array := make([]int, len(nodes))
	for i, s := range nodes {
		switch s {
		case '.':
			array[i] = int(stone.None)
		case 'x':
			array[i] = int(stone.Black)
		case 'o':
			array[i] = int(stone.White)
		case '+':
			array[i] = int(stone.Wall)
		default:
			// Ignored.
		}
	}

	return array
}

// BoardSize - 何路盤か。
func (config EngineConf) BoardSize() int {
	return int(config.Engine.BoardSize)
}

// SentinelBoardMax - 枠付きの盤上の交点の数
func (config EngineConf) SentinelBoardMax() int {
	// Width - 枠込み。
	Width := int(config.Engine.BoardSize) + 2
	// BoardMax - 枠込み盤の配列サイズ。
	return Width * Width
}

// Komi - float 32bit で足りるが、実行速度優先で float 64bit に変換して返します。
func (config EngineConf) Komi() float64 {
	return float64(config.Engine.Komi)
}

// MaxMoves - 最大手数。
func (config EngineConf) MaxMoves() int {
	return int(config.Engine.MaxMoves)
}
