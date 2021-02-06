package main

import (
	"fmt"
	"math/rand"
	"time"

	c "github.com/muzudho/kifuwarabe-uec12/controller"
	e "github.com/muzudho/kifuwarabe-uec12/entities"
	"github.com/ziutek/telnet"
)

// KifuwarabeV1 - きふわらべバージョン１。
// NNGSへの接続を試みる。
func KifuwarabeV1() {
	e.G.Chat.Trace("# きふわらべv1プログラム開始☆（＾～＾）\n")

	config := c.LoadGameConf("input/kifuwarabe-v1.gameConf.toml")

	/*
		e.G.Chat.Trace("# Config読んだ☆（＾～＾）\n")
		e.G.Chat.Trace("# Komi=%f\n", config.Game.Komi)
		e.G.Chat.Trace("# BoardSize=%d\n", config.Game.BoardSize)
		e.G.Chat.Trace("# MaxMoves=%d\n", config.Game.MaxMoves)
		e.G.Chat.Trace("# BoardData=%s\n", config.Game.BoardData)
		e.G.Chat.Trace("# SentinelBoardMax()=%d\n", config.SentinelBoardMax())
	*/

	board := e.NewBoardV9a(config.GetBoardArray(), config.BoardSize(), config.SentinelBoardMax(), config.Komi(), config.MaxMoves())
	// presenter := p.NewPresenterV9a()

	//e.G.Chat.Trace("# 盤を新規作成した☆（＾～＾）\n")

	rand.Seed(time.Now().UnixNano())

	//e.G.Chat.Trace("# (^q^) ランダムの種を設定したぜ☆\n")

	board.InitBoard()

	e.G.Chat.Trace("# NNGSへの接続を試みるぜ☆（＾～＾） server=%s port=%d\n", config.Nngs.Server, config.Nngs.Port)

	// connectionString := fmt.Sprintf("%s:%d", config.Nngs.Server, config.Nngs.Port)
	// connectionString := fmt.Sprintf("localhost:5555", config.Nngs.Server, config.Nngs.Port)

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
	e.G.Chat.Trace("# NNGSへ接続でけた☆（＾～＾）\n")

	e.G.Chat.Trace("# NNGSへユーザー名 %s を送ったろ……☆（＾～＾）\n", config.Nngs.User)

	nngsConn.Write([]byte(fmt.Sprintf("%s\n", config.Nngs.User)))

	e.G.Chat.Trace("# NNGSからの返信を待と……☆（＾～＾）\n")

	// nngsConnBuf := bufio.NewReader(nngsConn)
	// str, err := nngsConnBuf.ReadString('\n')

	// str, err := nngsConn.ReadUntil("\n")
	str, err := nngsConn.ReadString('\n')
	e.G.Chat.Trace("# どうか☆（＾～＾）\n")
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

	e.G.Chat.Trace("# NNGSへの接続終わった☆（＾～＾）\n")
}
