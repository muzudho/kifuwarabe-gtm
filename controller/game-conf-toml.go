package controller

import (
	"io/ioutil"
	"strconv"
	"strings"

	e "github.com/muzudho/kifuwarabe-uec12/entities"
	"github.com/pelletier/go-toml"
)

// Config - Tomlファイル。
type Config struct {
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
func (config Config) GetBoardArray() []int {
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
func (config Config) BoardSize() int {
	return int(config.Game.BoardSize)
}

// SentinelBoardMax - 枠付きの盤上の交点の数
func (config Config) SentinelBoardMax() int {
	// Width - 枠込み。
	Width := int(config.Game.BoardSize) + 2
	// BoardMax - 枠込み盤の配列サイズ。
	return Width * Width
}

// Komi - float 32bit で足りるが、実行速度優先で float 64bit に変換して返します。
func (config Config) Komi() float64 {
	return float64(config.Game.Komi)
}

// MaxMoves - 最大手数。
func (config Config) MaxMoves() int {
	return int(config.Game.MaxMoves)
}

// LoadGameConf - ゲーム設定ファイルを読み込みます。
func LoadGameConf(path string) Config {

	// ファイル読込
	fileData, err := ioutil.ReadFile(path)
	if err != nil {
		e.G.Chat.Fatal("path=%s", path)
		panic(err)
	}
	// fmt.Print(string(fileData))

	/*
		// Toml解析
		tomlTree, err := toml.Load(string(fileData))
		if err != nil {
			panic(err)
		}

		fmt.Println("Success.")

		komi := tomlTree.Get("Game.Komi").(float64)
		fmt.Printf("komi=%f\n", komi)

		boardSize := tomlTree.Get("Game.BoardSize").(int64)
		fmt.Printf("boardSize=%d\n", boardSize)

		maxMoves := tomlTree.Get("Game.MaxMoves").(int64)
		fmt.Printf("maxMoves=%d\n", maxMoves)
	*/

	// Toml解析
	binary := []byte(string(fileData))
	config := Config{}
	toml.Unmarshal(binary, &config)
	/*
		fmt.Println("Komi=", config.Game.Komi)
		fmt.Println("BoardSize=", config.Game.BoardSize)
		fmt.Println("MaxMoves=", config.Game.MaxMoves)
	*/

	return config
}
