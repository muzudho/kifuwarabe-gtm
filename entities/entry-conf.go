package entities

import (
	"strconv"
	"strings"
)

// EntryConf - Tomlファイル。
type EntryConf struct {
	Server Server
	User   User
	Engine Engine
}

// Server - [Server] テーブル。
type Server struct {
	Host string
	Port uint16
}

// User - [User] 区画。
type User struct {
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

// GetBoardArray - 盤上の石の色の配列。
func (config EntryConf) GetBoardArray() []int {
	// 最後のカンマを削除しないと、要素数が 1 多くなってしまいます。
	s := strings.TrimRight(config.Engine.BoardData, ",")
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
	return int(config.Engine.BoardSize)
}

// SentinelBoardMax - 枠付きの盤上の交点の数
func (config EntryConf) SentinelBoardMax() int {
	// Width - 枠込み。
	Width := int(config.Engine.BoardSize) + 2
	// BoardMax - 枠込み盤の配列サイズ。
	return Width * Width
}

// Komi - float 32bit で足りるが、実行速度優先で float 64bit に変換して返します。
func (config EntryConf) Komi() float64 {
	return float64(config.Engine.Komi)
}

// MaxMoves - 最大手数。
func (config EntryConf) MaxMoves() int {
	return int(config.Engine.MaxMoves)
}
