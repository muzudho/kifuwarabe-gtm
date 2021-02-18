package ui

import (
	"fmt"
	"io/ioutil"

	e "github.com/muzudho/kifuwarabe-gtp/entities"
	"github.com/pelletier/go-toml"
)

// LoadEntryConf - ゲーム設定ファイルを読み込みます。
func LoadEntryConf(path string) e.EntryConf {

	// ファイル読込
	fileData, err := ioutil.ReadFile(path)
	if err != nil {
		e.G.Log.Fatal("path=%s", path)
		panic(err)
	}

	debugPrintToml(fileData)

	// Toml解析
	binary := []byte(string(fileData))
	config := e.EntryConf{}
	toml.Unmarshal(binary, &config)

	return config
}

func debugPrintToml(fileData []byte) {
	fmt.Print(string(fileData))

	// Toml解析
	tomlTree, err := toml.Load(string(fileData))
	if err != nil {
		panic(err)
	}
	fmt.Println("[情報] Input:")
	fmt.Printf("Game.Komi=%f\n", tomlTree.Get("Game.Komi").(float64))
	fmt.Printf("Game.BoardSize=%d\n", tomlTree.Get("").(int64))
	fmt.Printf("Game.MaxMoves=%d\n", tomlTree.Get("").(int64))
}
func debugPrintConfig(config e.EntryConf) {
	fmt.Println("[情報] Memory:")
	fmt.Printf("[情報] Game.Komi=%f\n", config.Game.Komi)
	fmt.Printf("[情報] Game.Komi=%d\n", config.Game.BoardSize)
	fmt.Printf("[情報] Game.Komi=%d\n", config.Game.MaxMoves)
}
