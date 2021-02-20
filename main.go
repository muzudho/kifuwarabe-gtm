// Source: https://github.com/bleu48/GoGo
// 電通大で行われたコンピュータ囲碁講習会をGolangで追う

package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	e "github.com/muzudho/kifuwarabe-gtp/entities"
	p "github.com/muzudho/kifuwarabe-gtp/presenter"
	"github.com/muzudho/kifuwarabe-gtp/ui"
	u "github.com/muzudho/kifuwarabe-gtp/usecases"
)

func main() {
	// Working directory
	wdir, err := os.Getwd()
	if err != nil {
		// ここでは、ログはまだ設定できてない
		panic(fmt.Sprintf("<Engine> wdir=%s", wdir))
	}

	// コマンドライン引数
	workdir := flag.String("workdir", wdir, "Working directory path.")
	flag.Parse()
	entryConfPath := filepath.Join(*workdir, "input/default.entryConf.toml")

	// グローバル変数の作成
	u.G = *new(u.GlobalVariables)

	// ロガーの作成。
	// TODO ディレクトリが存在しなければ、強制終了します。
	u.G.Log = *u.NewLogger(
		filepath.Join(*workdir, "output/trace.log"),
		filepath.Join(*workdir, "output/debug.log"),
		filepath.Join(*workdir, "output/info.log"),
		filepath.Join(*workdir, "output/notice.log"),
		filepath.Join(*workdir, "output/warn.log"),
		filepath.Join(*workdir, "output/error.log"),
		filepath.Join(*workdir, "output/fatal.log"),
		filepath.Join(*workdir, "output/print.log"))

	u.G.Log.Trace("<Engine> KifuwarabeGoGo プログラム開始☆（＾～＾）\n")
	u.G.Log.Trace("<Engine> Author: %s\n", u.Author)
	u.G.Log.Trace("<Engine> This is a GTP engine.\n")
	u.G.Log.Trace("<Engine> wdir=%s\n", wdir)
	u.G.Log.Trace("<Engine> flag.Args()=%s\n", flag.Args())
	u.G.Log.Trace("<Engine> workdir=%s\n", *workdir)
	u.G.Log.Trace("<Engine> entryConfPath=%s\n", entryConfPath)

	// チャッターの作成。 標準出力とロガーを一緒にしただけです。
	u.G.Chat = *u.NewChatter(u.G.Log)

	// TODO ファイルが存在しなければ、強制終了します。
	config := ui.LoadEntryConf(entryConfPath) // "input/default.entryConf.toml"

	board := e.NewBoard(config.GetBoardArray(), config.BoardSize(), config.SentinelBoardMax(), config.Komi(), config.MaxMoves())
	e.UctChildrenSize = config.BoardSize()*config.BoardSize() + 1

	u.G.Log.Trace("<Engine> board.BoardSize()=%d\n", board.BoardSize())
	u.G.Log.Trace("<Engine> board.SentinelBoardMax()=%d\n", board.SentinelBoardMax())

	presenter := p.NewPresenterV9a()

	rand.Seed(time.Now().UnixNano())
	board.InitBoard()

	u.G.Log.Trace("<Engine> 何か標準入力しろだぜ☆（＾～＾）\n")

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		command := scanner.Text()
		u.G.Log.Trace("<Engine> '%s' コマンドを受け取ったぜ☆（＾～＾）\n", command)

		tokens := strings.Split(command, " ")
		switch tokens[0] {
		case "boardsize":
			u.G.Log.Trace("<Engine> 'boardsize' には空行を返すぜ（＾～＾）\n")
			u.G.Chat.Print("= \n\n")
		case "clear_board":
			board.InitBoard()
			u.G.Log.Trace("<Engine> 'clear_board' に対応して盤を初期化するぜ（＾～＾）\n")
			u.G.Chat.Print("= \n\n")
		case "quit":
			u.G.Log.Trace("<Engine> 'quit' に対応してアプリケーションを終了するぜ（＾～＾）\n")
			os.Exit(0)
		case "protocol_version":
			u.G.Log.Trace("<Engine> 'protocol_version' に対応してバージョン番号を返すぜ（＾～＾）\n")
			u.G.Chat.Print("= 2\n\n")
		case "name":
			u.G.Log.Trace("<Engine> 'name' に対応して名前を返すぜ（＾～＾）\n")
			u.G.Chat.Print("= KwGoGo\n\n")
		case "version":
			u.G.Log.Trace("<Engine> 'version' に対応してバージョン番号を返すぜ（＾～＾）\n")
			u.G.Chat.Print("= 0.0.1\n\n")
		case "list_commands":
			u.G.Log.Trace("<Engine> 'list_commands' に対応してコマンドのリストを返すぜ（＾～＾）\n")
			u.G.Chat.Print("= boardsize\nclear_board\nquit\nprotocol_version\nundo\n" +
				"name\nversion\nlist_commands\nkomi\ngenmove\nplay\n\n")
		case "komi":
			u.G.Log.Trace("<Engine> 'komi' に対応してコミを返すぜ（＾～＾）\n")
			u.G.Chat.Print("= 6.5\n\n") // TODO コミ
		case "undo":
			u.UndoV9()
			u.G.Log.Trace("<Engine> 'undo' は未実装だぜ（＾～＾）\n")
			u.G.Chat.Print("= \n\n")
		// 19路盤だと、すごい長い時間かかる。
		// genmove b
		case "genmove":
			color := 1
			if 1 < len(tokens) && strings.ToLower(tokens[1]) == "w" {
				color = 2
			}
			tIdx := u.PlayComputerMove(board, color, 1, presenter.PrintBoardType1, presenter.PrintBoardType2)

			bestmoveString := p.GetPointName(board, tIdx)

			u.G.Log.Trace("<Engine> 'genmove' に対応して[%s]を返すぜ（＾～＾）\n", bestmoveString)
			u.G.Chat.Print("= %s\n\n", bestmoveString)
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
			// u.G.Log.Trace("<Engine> 'play' コマンドを受け取ったぜ☆（＾～＾）\n")

			color := 1
			if 1 < len(tokens) && strings.ToLower(tokens[1]) == "w" {
				color = 2
			}

			u.G.Log.Trace("<Engine> color=%d len(tokens)=%d\n", color, len(tokens))

			if 2 < len(tokens) {
				u.G.Log.Trace("<Engine> tokens[2]=%s\n", tokens[2])
				var tIdx int
				if strings.ToLower(tokens[2]) == "pass" {
					tIdx = 0
					u.G.Log.Trace("<Engine> Pass\n")
				} else {
					x, y, err := e.GetXYFromName(tokens[2])
					if err != nil {
						panic(err)
					}

					tIdx = board.GetTIdxFromXy(x, y)

					u.G.Log.Trace("<Engine> x=%d y=%d\n", x, y)
				}
				board.AddMoves(tIdx, color, 0, presenter.PrintBoardType2)

				u.G.Log.Trace("<Engine> 'play' に対応して空行を返すぜ（＾～＾）\n")
				u.G.Chat.Print("= \n\n")
			}
		default:
			u.G.Log.Trace("<Engine> '%s' コマンドには未対応だぜ（＾～＾）\n", tokens[0])
			u.G.Chat.Print("? unknown_command\n\n")
		}
	}
}
