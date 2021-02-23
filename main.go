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
	"github.com/muzudho/kifuwarabe-gtp/presenter"
	p "github.com/muzudho/kifuwarabe-gtp/presenter"
	"github.com/muzudho/kifuwarabe-gtp/ui"
	u "github.com/muzudho/kifuwarabe-gtp/usecases"
)

func main() {
	// Working directory
	wdir, err := os.Getwd()
	if err != nil {
		// ここでは、ログはまだ設定できてない
		panic(fmt.Sprintf("...Engine wdir=%s", wdir))
	}

	// コマンドライン引数
	workdir := flag.String("workdir", wdir, "Working directory path.")
	flag.Parse()
	engineConfPath := filepath.Join(*workdir, "input/engine.conf.toml")

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

	// 既存のログ・ファイルを削除
	u.G.Log.RemoveAllOldLogs()

	// ログ・ファイルの開閉
	err = u.G.Log.OpenAllLogs()
	if err != nil {
		panic(u.G.Log.Fatal(err.Error()))
	}
	defer u.G.Log.CloseAllLogs()

	u.G.Log.Trace("...Engine Remove all old logs\n")
	u.G.Log.Trace("...Engine KifuwarabeGoGo プログラム開始☆（＾～＾）\n")
	u.G.Log.Trace("...Engine Author: %s\n", u.Author)
	u.G.Log.Trace("...Engine This is a GTP engine.\n")
	u.G.Log.Trace("...Engine wdir=%s\n", wdir)
	u.G.Log.Trace("...Engine flag.Args()=%s\n", flag.Args())
	u.G.Log.Trace("...Engine workdir=%s\n", *workdir)
	u.G.Log.Trace("...Engine engineConfPath=%s\n", engineConfPath)

	// チャッターの作成。 標準出力とロガーを一緒にしただけです。
	u.G.Chat = *u.NewChatter(u.G.Log)
	u.G.StderrChat = *u.NewStderrChatter(u.G.Log)

	// TODO ファイルが存在しなければ、強制終了します。
	config, err := ui.LoadEngineConf(engineConfPath)
	if err != nil {
		panic(u.G.Log.Fatal("path=[%s] err=[%s]", engineConfPath, err))
	}

	rand.Seed(time.Now().UnixNano())

	position := e.NewPosition(config.GetBoardArray(), config.BoardSize(), config.SentinelBoardMax(), config.Komi(), config.MaxMoves())
	u.G.Log.Trace("...Engine position.BoardSize()=%d\n", position.BoardSize())
	u.G.Log.Trace("...Engine position.SentinelBoardMax()=%d\n", position.SentinelBoardMax())
	e.UctChildrenSize = config.BoardSize()*config.BoardSize() + 1

	u.G.Log.Trace("...Engine 何か標準入力しろだぜ☆（＾～＾）\n")

	scanner := bufio.NewScanner(os.Stdin)

MailLoop:
	for scanner.Scan() {
		u.G.Log.FlushAllLogs()

		command := scanner.Text()
		u.G.Log.Notice("-->%s '%s' command\n", config.Profile.Name, command)

		tokens := strings.Split(command, " ")
		switch tokens[0] {
		case "boardsize":
			u.G.Log.Notice("<--%s ok\n", config.Profile.Name)
			u.G.Chat.Print("= \n\n")
		case "clear_board":
			position = e.NewPosition(config.GetBoardArray(), config.BoardSize(), config.SentinelBoardMax(), config.Komi(), config.MaxMoves())
			u.G.Log.Notice("<--%s ok\n", config.Profile.Name)
			u.G.Chat.Print("= \n\n")
		case "quit":
			u.G.Log.Notice("<--%s Quit\n", config.Profile.Name)
			break MailLoop
			// os.Exit(0)
		case "protocol_version":
			u.G.Log.Notice("<--%s Version ok\n", config.Profile.Name)
			u.G.Chat.Print("= 2\n\n")
		case "name":
			u.G.Log.Notice("<--%s Name ok\n", config.Profile.Name)
			u.G.Chat.Print("= KwGoGo\n\n")
		case "version":
			u.G.Log.Notice("<--%s Version ok\n", config.Profile.Name)
			u.G.Chat.Print("= 0.0.1\n\n")
		case "list_commands":
			u.G.Log.Notice("<--%s CommandList ok\n", config.Profile.Name)
			u.G.Chat.Print("= boardsize\nclear_board\nquit\nprotocol_version\nundo\n" +
				"name\nversion\nlist_commands\nkomi\ngenmove\nplay\n\n")
		case "komi":
			u.G.Log.Notice("<--%s Komi ok\n", config.Profile.Name)
			u.G.Chat.Print("= 6.5\n\n") // TODO コミ
		case "undo":
			u.UndoV9() // TODO アンドゥ
			u.G.Log.Notice("<--%s Unimplemented undo, ignored\n", config.Profile.Name)
			u.G.Chat.Print("= \n\n")
		// 19路盤だと、すごい長い時間かかる。
		// genmove b
		case "genmove":
			color := 1
			if 1 < len(tokens) && strings.ToLower(tokens[1]) == "w" {
				color = 2
			}
			tIdx := u.PlayComputerMove(position, color, 1, presenter.PrintBoard)
			presenter.PrintBoardHeader(position, position.MovesNum)
			presenter.PrintBoard(position)

			bestmoveString := p.GetPointName(position, tIdx)

			u.G.Log.Notice("<--%s [%s] ok\n", config.Profile.Name, bestmoveString)
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
			color := 1
			if 1 < len(tokens) && strings.ToLower(tokens[1]) == "w" {
				color = 2
			}

			// u.G.Log.Trace("...Engine color=%d len(tokens)=%d\n", color, len(tokens))

			if 2 < len(tokens) {
				// u.G.Log.Trace("...Engine tokens[2]=%s\n", tokens[2])
				var tIdx int
				if strings.ToLower(tokens[2]) == "pass" {
					tIdx = 0
					// u.G.Log.Trace("...Engine pass\n")
				} else {
					x, y, err := e.GetXYFromName(tokens[2])
					if err != nil {
						panic(u.G.Log.Fatal(err.Error()))
					}

					tIdx = position.GetTIdxFromFileRank(x+1, y+1)

					// u.G.Log.Trace("...Engine file=%d rank=%d\n", x+1, y+1)
				}
				position.AddMoves(tIdx, color, 0)
				presenter.PrintBoardHeader(position, position.MovesNum)
				presenter.PrintBoard(position)

				u.G.Log.Notice("<--%s ok\n", config.Profile.Name)
				u.G.Chat.Print("= \n\n")
			}
		default:
			u.G.Log.Notice("<--%s Unimplemented '%s' command\n", config.Profile.Name, tokens[0])
			u.G.Chat.Print("? unknown_command\n\n")
		}
	}

	u.G.Log.Trace("...%s... End engine\n", config.Profile.Name)
}
