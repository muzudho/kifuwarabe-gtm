package main

import (
	"math/rand"
	"time"

	e "github.com/muzudho/kifuwarabe-gtp/entities"
	"github.com/muzudho/kifuwarabe-gtp/ui"
)

// KifuwarabeV1 - きふわらべバージョン１。
// NNGSへの接続を試みる。
func KifuwarabeV1() {
	e.G.Log.Trace("# きふわらべv1プログラム開始☆（＾～＾）\n")

	config := ui.LoadEntryConf("input/kifuwarabe-v1.entryConf.toml")

	board := e.NewBoard(config.GetBoardArray(), config.BoardSize(), config.SentinelBoardMax(), config.Komi(), config.MaxMoves())
	// presenter := p.NewPresenterV9a()

	//e.G.Log.Trace("# 盤を新規作成した☆（＾～＾）\n")

	rand.Seed(time.Now().UnixNano())

	//e.G.Log.Trace("# (^q^) ランダムの種を設定したぜ☆\n")

	board.InitBoard()
}
