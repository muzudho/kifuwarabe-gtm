package main

import (
	"fmt"
	"math/rand"
	"time"

	e "github.com/muzudho/kifuwarabe-gtp/entities"
	"github.com/muzudho/kifuwarabe-gtp/ui"
	"github.com/ziutek/telnet"
)

// KifuwarabeV1 - きふわらべバージョン１。
// NNGSへの接続を試みる。
func KifuwarabeV1() {
	e.G.Log.Trace("# きふわらべv1プログラム開始☆（＾～＾）\n")

	config := ui.LoadEntryConf("input/kifuwarabe-v1.entryConf.toml")

	/*
		e.G.Log.Trace("# Config読んだ☆（＾～＾）\n")
		e.G.Log.Trace("# Komi=%f\n", config.Game.Komi)
		e.G.Log.Trace("# BoardSize=%d\n", config.Game.BoardSize)
		e.G.Log.Trace("# MaxMoves=%d\n", config.Game.MaxMoves)
		e.G.Log.Trace("# BoardData=%s\n", config.Game.BoardData)
		e.G.Log.Trace("# SentinelBoardMax()=%d\n", config.SentinelBoardMax())
	*/

	board := e.NewBoard(config.GetBoardArray(), config.BoardSize(), config.SentinelBoardMax(), config.Komi(), config.MaxMoves())
	// presenter := p.NewPresenterV9a()

	//e.G.Log.Trace("# 盤を新規作成した☆（＾～＾）\n")

	rand.Seed(time.Now().UnixNano())

	//e.G.Log.Trace("# (^q^) ランダムの種を設定したぜ☆\n")

	board.InitBoard()

	e.G.Log.Trace("# NNGSへの接続を試みるぜ☆（＾～＾） host=%s port=%d\n", config.Server.Host, config.Server.Port)

	// connectionString := fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)
	// connectionString := fmt.Sprintf("localhost:5555", config.Server.Host, config.Server.Port)

	// "tcp" で合ってるよう。
	nngsConn, err := telnet.Dial("tcp", "localhost:5555")
	// nngsConn, err := telnet.Dial("udp", "localhost:5555")
	// fail: nngsConn, err := telnet.Dial("ip4", "localhost:5555")
	// fail: nngsConn, err := telnet.Dial("ip", "localhost:5555")
	// nngsConn, err := telnet.Dial("tcp", connectionString)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect. %s", err))
	}
	defer nngsConn.Close()
	e.G.Log.Trace("# NNGSへ接続でけた☆（＾～＾）\n")

	e.G.Log.Trace("# NNGSへユーザー名 %s を送ったろ……☆（＾～＾）\n", config.User.Name)

	nngsConn.Write([]byte(fmt.Sprintf("%s\n", config.User.Name)))

	e.G.Log.Trace("# NNGSからの返信を待と……☆（＾～＾）\n")

	// nngsConnBuf := bufio.NewReader(nngsConn)
	// str, err := nngsConnBuf.ReadString('\n')

	// str, err := nngsConn.ReadUntil("\n")
	str, err := nngsConn.ReadString('\n')
	e.G.Log.Trace("# どうか☆（＾～＾）\n")
	if err != nil {
		panic(fmt.Sprintf("Failed to read string. %s", err))
	}
	fmt.Printf("str=%s", str)

	/*
		// scanner - 標準入力を監視します。
		scanner := bufio.NewScanner(os.Stdin)
		// 一行読み取ります。
		for scanner.Scan() {
			// 書き込みます。最後に改行を付けます。
			oi.LongWrite(w, scanner.Bytes())
			oi.LongWrite(w, []byte("\n"))
		}
	*/

	e.G.Log.Trace("# NNGSへの接続終わった☆（＾～＾）\n")
}
