package presenter

import (
	"fmt"
	"strings"

	e "github.com/muzudho/kifuwarabe-gtp/entities"
	u "github.com/muzudho/kifuwarabe-gtp/usecases"
)

// BoardView - 表示機能 Version 9a.
type BoardView struct {
}

// labelOfColumns - 各列の表示符号。
// I は欠番です。
var labelOfColumns = [20]byte{'@', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'J',
	'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T'}

// labelOfRows - 各行の表示符号。
var labelOfRows = [20]string{" 0", " 1", " 2", " 3", " 4", " 5", " 6", " 7", " 8", " 9",
	"10", "11", "12", "13", "14", "15", "16", "17", "18", "19"}

// " x" - Visual Studio Code の 全角半角崩れ対応。
// " ○" - Visual Studio Code の 全角半角崩れ対応。
var stoneLabels = [4]string{" .", " x", " o", " #"}

// PrintBoardHeader - 手数などを表示
func PrintBoardHeader(board *e.Board, moves int) {
	u.G.StderrChat.Info("[ Ko=%s MovesNum=%d ]\n", (*board).GetNameFromTIdx(board.KoIdx), moves)
}

// PrintBoard - 盤を描画
func PrintBoard(board *e.Board) {
	boardSize := (*board).BoardSize()

	var b strings.Builder
	b.Grow(3 * boardSize) // だいたい適当

	b.WriteString("\n   ")
	for x := 0; x < boardSize; x++ {
		b.WriteString(fmt.Sprintf(" %c", labelOfColumns[x+1]))
	}
	b.WriteString("\n  +")
	for x := 0; x < boardSize; x++ {
		b.WriteString("--")
	}
	b.WriteString("+\n")
	for y := 0; y < boardSize; y++ {
		b.WriteString(fmt.Sprintf("%s|", labelOfRows[y+1]))
		for x := 0; x < boardSize; x++ {
			b.WriteString(fmt.Sprintf("%s", stoneLabels[(*board).ColorAtFileRank(x+1, y+1)]))
		}
		b.WriteString("|\n")
	}
	b.WriteString("  +")
	for x := 0; x < boardSize; x++ {
		b.WriteString("--")
	}
	b.WriteString("+\n")

	u.G.StderrChat.Info(b.String())
}

// PrintSgf - SGF形式の棋譜表示
func PrintSgf(board *e.Board, moves int, record []int) {
	boardSize := board.BoardSize()

	fmt.Printf("(;GM[1]SZ[%d]KM[%.1f]PB[]PW[]\n", boardSize, board.Komi())
	for i := 0; i < moves; i++ {
		z := record[i]
		y := z / board.SentinelWidth()
		x := z - y*board.SentinelWidth()
		var sStone = [2]string{"B", "W"}
		fmt.Printf(";%s", sStone[i&1])
		if z == 0 {
			fmt.Printf("[]")
		} else {
			fmt.Printf("[%c%c]", x+'a'-1, y+'a'-1)
		}
		if ((i + 1) % 10) == 0 {
			fmt.Printf("\n")
		}
	}
	fmt.Printf(")\n")
}

// GetPointName - YX座標の文字表示？ A1 とか
func GetPointName(board *e.Board, tIdx int) string {
	if tIdx == 0 {
		return "pass"
	}

	// boardSize := board.BoardSize()

	y := tIdx / (*board).SentinelWidth()
	x := tIdx - y*(*board).SentinelWidth()

	ax := labelOfColumns[x]

	//return string(ax) + string(boardSize+1-y+'0')
	return fmt.Sprintf("%c%d", ax, y)
}
