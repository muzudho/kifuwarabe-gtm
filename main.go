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

	l "github.com/muzudho/go-logger"
	be "github.com/muzudho/kifuwarabe-go-base/entities"
	tbe "github.com/muzudho/kifuwarabe-go-think-base/entities"
	g "github.com/muzudho/kifuwarabe-gtp/global"
	"github.com/muzudho/kifuwarabe-gtp/presenter"
	p "github.com/muzudho/kifuwarabe-gtp/presenter"
	"github.com/muzudho/kifuwarabe-gtp/ui"
	u "github.com/muzudho/kifuwarabe-gtp/usecases"
)

func main() {
	// Working directory
	dwd, err := os.Getwd()
	if err != nil {
		// ここでは、ログはまだ設定できてない
		panic(fmt.Sprintf("...Engine DefaultWorkingDirectory=%s", dwd))
	}

	// コマンドライン引数登録
	workdir := flag.String("workdir", dwd, "Working directory path.")
	// 解析
	flag.Parse()

	engineConfPath := filepath.Join(*workdir, "input/engine.conf.toml")

	tracePath := filepath.Join(*workdir, "output/trace.log")
	debugPath := filepath.Join(*workdir, "output/debug.log")
	infoPath := filepath.Join(*workdir, "output/info.log")
	noticePath := filepath.Join(*workdir, "output/notice.log")
	warnPath := filepath.Join(*workdir, "output/warn.log")
	errorPath := filepath.Join(*workdir, "output/error.log")
	fatalPath := filepath.Join(*workdir, "output/fatal.log")
	printPath := filepath.Join(*workdir, "output/print.log")
	/* Debug
	fmt.Printf("...Engine tracePath=%s\n", tracePath)
	fmt.Printf("...Engine debugPath=%s\n", debugPath)
	fmt.Printf("...Engine infoPath=%s\n", infoPath)
	fmt.Printf("...Engine noticePath=%s\n", noticePath)
	fmt.Printf("...Engine warnPath=%s\n", warnPath)
	fmt.Printf("...Engine errorPath=%s\n", errorPath)
	fmt.Printf("...Engine fatalPath=%s\n", fatalPath)
	fmt.Printf("...Engine printPath=%s\n", printPath)
	*/

	// グローバル変数の作成
	g.G = *new(g.Variables)

	// ロガーの作成。
	// TODO ディレクトリが存在しなければ、強制終了します。
	g.G.Log = *l.NewLogger(
		tracePath,
		debugPath,
		infoPath,
		noticePath,
		warnPath,
		errorPath,
		fatalPath,
		printPath)

	// 既存のログ・ファイルを削除
	g.G.Log.RemoveAllOldLogs()

	// ログ・ファイルの開閉
	err = g.G.Log.OpenAllLogs()
	if err != nil {
		// ログ・ファイルを開くのに失敗したのだから、ログ・ファイルへは書き込めません
		panic(fmt.Sprintf("...Engine... %s", err))
	}
	defer g.G.Log.CloseAllLogs()

	g.G.Log.Trace("...Engine Remove all old logs\n")
	g.G.Log.Trace("...Engine KifuwarabeGoGo プログラム開始☆（＾～＾）\n")
	g.G.Log.Trace("...Engine Author: %s\n", g.Author)
	g.G.Log.Trace("...Engine This is a GTP engine.\n")
	g.G.Log.Trace("...Engine DefaultWorkingDirectory=%s\n", dwd)
	g.G.Log.Trace("...Engine flag.Args()=%s\n", flag.Args())
	g.G.Log.Trace("...Engine workdir=%s\n", *workdir)
	g.G.Log.Trace("...Engine engineConfPath=%s\n", engineConfPath)

	// チャッターの作成。 標準出力とロガーを一緒にしただけです。
	g.G.Chat = *l.NewChatter(g.G.Log)
	g.G.StderrChat = *l.NewStderrChatter(g.G.Log)

	// TODO ファイルが存在しなければ、強制終了します。
	config, err := ui.LoadEngineConf(engineConfPath)
	if err != nil {
		panic(g.G.Log.Fatal(fmt.Sprintf("...Engine... path=[%s] err=[%s]", engineConfPath, err)))
	}

	rand.Seed(time.Now().UnixNano())

	position := be.NewPosition(config.GetBoardArray(), config.BoardSize(), config.SentinelBoardMax(), config.Komi(), config.MaxMoves())
	g.G.Log.Trace("...Engine position.BoardSize()=%d\n", position.BoardSize())
	g.G.Log.Trace("...Engine position.SentinelBoardMax()=%d\n", position.SentinelBoardMax())
	tbe.UctChildrenSize = config.BoardSize()*config.BoardSize() + 1

	g.G.Log.Trace("...Engine 何か標準入力しろだぜ☆（＾～＾）\n")

	scanner := bufio.NewScanner(os.Stdin)

MainLoop:
	for scanner.Scan() {
		g.G.Log.FlushAllLogs()

		command := scanner.Text()
		g.G.Log.Notice("-->%s '%s' command\n", config.Profile.Name, command)

		tokens := strings.Split(command, " ")
		switch tokens[0] {
		case "boardsize":
			g.G.Log.Notice("<--%s ok\n", config.Profile.Name)
			g.G.Chat.Print("= \n\n")
		case "clear_board":
			position = be.NewPosition(config.GetBoardArray(), config.BoardSize(), config.SentinelBoardMax(), config.Komi(), config.MaxMoves())
			g.G.Log.Notice("<--%s ok\n", config.Profile.Name)
			g.G.Chat.Print("= \n\n")
		case "quit":
			g.G.Log.Notice("<--%s Quit\n", config.Profile.Name)
			break MainLoop
			// os.Exit(0)
		case "protocol_version":
			g.G.Log.Notice("<--%s Version ok\n", config.Profile.Name)
			g.G.Chat.Print("= 2\n\n")
		case "name":
			g.G.Log.Notice("<--%s Name ok\n", config.Profile.Name)
			g.G.Chat.Print("= KwGoGo\n\n")
		case "version":
			g.G.Log.Notice("<--%s Version ok\n", config.Profile.Name)
			g.G.Chat.Print("= 0.0.1\n\n")
		case "list_commands":
			g.G.Log.Notice("<--%s CommandList ok\n", config.Profile.Name)
			g.G.Chat.Print("= boardsize\nclear_board\nquit\nprotocol_version\nundo\n" +
				"name\nversion\nlist_commands\nkomi\ngenmove\nplay\n\n")
		case "komi":
			g.G.Log.Notice("<--%s Komi ok\n", config.Profile.Name)
			g.G.Chat.Print("= 6.5\n\n") // TODO コミ
		case "undo":
			u.UndoV9() // TODO アンドゥ
			g.G.Log.Notice("<--%s Unimplemented undo, ignored\n", config.Profile.Name)
			g.G.Chat.Print("= \n\n")
		// 19路盤だと、すごい長い時間かかる。
		// genmove b
		case "genmove":
			color := 1
			if 1 < len(tokens) && strings.ToLower(tokens[1]) == "w" {
				color = 2
			}
			tIdx := u.PlayComputerMove(position, color, 1, presenter.CreateBoardString)
			g.G.StderrChat.Info(presenter.CreateBoardHeader(position, position.MovesNum))
			g.G.StderrChat.Info(presenter.CreateBoardString(position))

			bestmoveString := p.GetPointName(position, tIdx)

			g.G.Log.Notice("<--%s [%s] ok\n", config.Profile.Name, bestmoveString)
			g.G.Chat.Print("= %s\n\n", bestmoveString)
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

			// g.G.Log.Trace("...Engine color=%d len(tokens)=%d\n", color, len(tokens))

			if 2 < len(tokens) {
				// g.G.Log.Trace("...Engine tokens[2]=%s\n", tokens[2])
				var tIdx int
				if strings.ToLower(tokens[2]) == "pass" {
					tIdx = 0
					// g.G.Log.Trace("...Engine pass\n")
				} else {
					x, y, err := be.GetXYFromName(tokens[2])
					if err != nil {
						panic(g.G.Log.Fatal(fmt.Sprintf("...Engine... %s", err)))
					}

					tIdx = position.GetTIdxFromFileRank(x+1, y+1)

					// g.G.Log.Trace("...Engine file=%d rank=%d\n", x+1, y+1)
				}
				position.AddMoves(tIdx, color, 0)
				g.G.StderrChat.Info(presenter.CreateBoardHeader(position, position.MovesNum))
				g.G.StderrChat.Info(presenter.CreateBoardString(position))

				g.G.Log.Notice("<--%s ok\n", config.Profile.Name)
				g.G.Chat.Print("= \n\n")
			}
		default:
			g.G.Log.Notice("<--%s Unimplemented '%s' command\n", config.Profile.Name, tokens[0])
			g.G.Chat.Print("? unknown_command\n\n")
		}
	}

	g.G.Log.Trace("...%s... End engine\n", config.Profile.Name)
	g.G.Log.FlushAllLogs()
}
