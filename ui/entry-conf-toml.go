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
	fmt.Printf("[情報] Engine.Komi=%f\n", tomlTree.Get("Engine.Komi").(float64))
	fmt.Printf("[情報] Engine.BoardSize=%d\n", tomlTree.Get("Engine.BoardSize").(int64))
	fmt.Printf("[情報] Engine.MaxMoves=%d\n", tomlTree.Get("Engine.MaxMoves").(int64))
	fmt.Printf("[情報] Engine.BoardData=%s\n", tomlTree.Get("Engine.BoardData").(string))
}
func debugPrintConfig(config e.EntryConf) {
	fmt.Println("[情報] Memory:")
	fmt.Printf("[情報] Engine.Komi=%f\n", config.Engine.Komi)
	fmt.Printf("[情報] Engine.BoardSize=%d\n", config.Engine.BoardSize)
	fmt.Printf("[情報] Engine.MaxMoves=%d\n", config.Engine.MaxMoves)
	fmt.Printf("[情報] Engine.MaxMoves=%s\n", config.Engine.BoardData)
}
