// Source: https://github.com/bleu48/GoGo
// 電通大で行われたコンピュータ囲碁講習会をGolangで追う

package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	c "github.com/muzudho/kifuwarabe-gtp/controller"
	e "github.com/muzudho/kifuwarabe-gtp/entities"
	p "github.com/muzudho/kifuwarabe-gtp/presenter"
	u "github.com/muzudho/kifuwarabe-gtp/usecases"
)

func main() {
	// グローバル変数の作成
	e.G = *new(e.GlobalVariables)

	// ロガーの作成。
	e.G.Log = *e.NewLogger(
		"output/trace.log",
		"output/debug.log",
		"output/info.log",
		"output/notice.log",
		"output/warn.log",
		"output/error.log",
		"output/fatal.log",
		"output/print.log")

	// チャッターの作成。 標準出力とロガーを一緒にしただけです。
	e.G.Chat = *e.NewChatter(e.G.Log)

	// 標準出力への表示と、ログへの書き込みを同時に行います。
	e.G.Log.Trace("Author: %s\n", e.Author)

	GoGoV9a() // GTP
	//KifuwarabeV1()
}

// GoGoV9a - バージョン９a。
// GTP2NNGS に対応しているのでは？
func GoGoV9a() {
	e.G.Log.Trace("# GoGo v9a プログラム開始☆（＾～＾）\n")

	config := c.LoadGameConf("input/example-v3.gameConf.toml")

	e.G.Log.Trace("# Config読んだ☆（＾～＾）\n")
	e.G.Log.Trace("# Server=%s\n", config.Nngs.Server)
	e.G.Log.Trace("# Port=%d\n", config.Nngs.Port)
	e.G.Log.Trace("# User=%s\n", config.Nngs.User)
	e.G.Log.Trace("# Pass=%s\n", config.Nngs.Pass)
	e.G.Log.Trace("# Komi=%f\n", config.Game.Komi)
	e.G.Log.Trace("# BoardSize=%d\n", config.Game.BoardSize)
	e.G.Log.Trace("# MaxMoves=%d\n", config.Game.MaxMoves)
	e.G.Log.Trace("# BoardData=%s\n", config.Game.BoardData)
	e.G.Log.Trace("# SentinelBoardMax()=%d\n", config.SentinelBoardMax())

	board := e.NewBoard(config.GetBoardArray(), config.BoardSize(), config.SentinelBoardMax(), config.Komi(), config.MaxMoves())

	e.G.Log.Trace("# board.BoardSize()=%d\n", board.BoardSize())
	e.G.Log.Trace("# board.SentinelBoardMax()=%d\n", board.SentinelBoardMax())
	// e.G.Log.Trace("# board.GetData()=", board.GetData())

	presenter := p.NewPresenterV9a()

	rand.Seed(time.Now().UnixNano())
	board.InitBoard()

	e.G.Log.Trace("何か標準入力しろだぜ☆（＾～＾）\n")

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		command := scanner.Text()
		tokens := strings.Split(command, " ")
		switch tokens[0] {
		case "boardsize":
			e.G.Chat.Print("= \n\n")
		case "clear_board":
			board.InitBoard()
			e.G.Chat.Print("= \n\n")
		case "quit":
			os.Exit(0)
		case "protocol_version":
			e.G.Chat.Print("= 2\n\n")
		case "name":
			e.G.Chat.Print("= GoGo\n\n")
		case "version":
			e.G.Chat.Print("= 0.0.1\n\n")
		case "list_commands":
			e.G.Chat.Print("= boardsize\nclear_board\nquit\nprotocol_version\nundo\n" +
				"name\nversion\nlist_commands\nkomi\ngenmove\nplay\n\n")
		case "komi":
			e.G.Chat.Print("= 6.5\n\n")
		case "undo":
			u.UndoV9()
			e.G.Chat.Print("= \n\n")
		// 19路盤だと、すごい長い時間かかる。
		// genmove b
		case "genmove":
			color := 1
			if 1 < len(tokens) && strings.ToLower(tokens[1]) == "w" {
				color = 2
			}
			z := u.PlayComputerMoveV9a(board, color, 1, presenter.PrintBoardType1, presenter.PrintBoardType2)
			e.G.Chat.Print("= %s\n\n", p.GetCharZ(board, z))
		// play b a3
		// play w d4
		// play b d5
		// play w e5
		// play b e4
		// play w d6
		// play b f5
		// play w c5
		// play b pass
		// play w pass
		case "play":
			color := 1
			if 1 < len(tokens) && strings.ToLower(tokens[1]) == "w" {
				color = 2
			}

			if 2 < len(tokens) {
				ax := strings.ToLower(tokens[2])
				fmt.Fprintf(os.Stderr, "ax=%s\n", ax)
				x := ax[0] - 'a' + 1
				if ax[0] >= 'i' {
					x--
				}
				y := int(ax[1] - '0')
				tIdx := board.GetTIdxFromXy(int(x)-1, board.BoardSize()-y)
				fmt.Fprintf(os.Stderr, "x=%d y=%d z=%04d\n", x, y, board.GetZ4(tIdx))
				if ax == "pass" {
					tIdx = 0
				}
				board.AddMovesType2(tIdx, color, 0, presenter.PrintBoardType2)
				e.G.Chat.Print("= \n\n")
			}
		default:
			e.G.Chat.Print("? unknown_command\n\n")
		}
	}
}
