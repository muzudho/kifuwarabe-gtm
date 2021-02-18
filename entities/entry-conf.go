package entities

import (
	"strconv"
	"strings"
)

// EntryConf - Tomlファイル。
type EntryConf struct {
	Nngs Nngs
	Game Game
}

// Nngs - [Nngs] テーブル。
type Nngs struct {
	Server string
	Port   uint16
	User   string
	Pass   string
}

// Game - [Game] テーブル。
type Game struct {
	Komi      float32
	BoardSize int8
	MaxMoves  int16
	BoardData string
}

// GetBoardArray - 盤上の石の色の配列。
func (config EntryConf) GetBoardArray() []int {
	// 最後のカンマを削除しないと、要素数が 1 多くなってしまいます。
	s := strings.TrimRight(config.Game.BoardData, ",")
	// fmt.Println("s=", s)
	nodes := strings.Split(s, ",")
	array := make([]int, len(nodes))
	for i, s := range nodes {
		s := strings.Trim(s, " ")
		color, _ := strconv.Atoi(s)
		// fmt.Println("strconv.Atoi(", s, ")=", color)
		array[i] = color
	}

	// fmt.Println("nodes=", nodes)
	// fmt.Println("array=", array)
	return array
}

// BoardSize - 何路盤か。
func (config EntryConf) BoardSize() int {
	return int(config.Game.BoardSize)
}

// SentinelBoardMax - 枠付きの盤上の交点の数
func (config EntryConf) SentinelBoardMax() int {
	// Width - 枠込み。
	Width := int(config.Game.BoardSize) + 2
	// BoardMax - 枠込み盤の配列サイズ。
	return Width * Width
}

// Komi - float 32bit で足りるが、実行速度優先で float 64bit に変換して返します。
func (config EntryConf) Komi() float64 {
	return float64(config.Game.Komi)
}

// MaxMoves - 最大手数。
func (config EntryConf) MaxMoves() int {
	return int(config.Game.MaxMoves)
}
