package presenter

import (
	"fmt"
	"os"

	e "github.com/muzudho/kifuwarabe-gtp/entities"
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

// PrintBoardType1 - 盤の描画。
func (presenter *BoardView) PrintBoardType1(board *e.Board) {
	boardSize := board.BoardSize()

	fmt.Printf("\n   ")
	for x := 0; x < boardSize; x++ {
		fmt.Printf(" %c", labelOfColumns[x+1])
	}
	fmt.Printf("\n  +")
	for x := 0; x < boardSize; x++ {
		fmt.Printf("--")
	}
	fmt.Printf("+\n")
	for y := 0; y < boardSize; y++ {
		fmt.Printf("%s|", labelOfRows[y+1])
		for x := 0; x < boardSize; x++ {
			fmt.Printf("%s", stoneLabels[board.ColorAtFileRank(x+1, y+1)])
		}
		fmt.Printf("|\n")
	}
	fmt.Printf("  +")
	for x := 0; x < boardSize; x++ {
		fmt.Printf("--")
	}
	fmt.Printf("+\n")
}

// PrintBoardType2 - 盤を描画。
func (presenter *BoardView) PrintBoardType2(board *e.Board, moves int) {
	boardSize := (*board).BoardSize()

	fmt.Fprintf(os.Stderr, "\n   ")
	for x := 0; x < boardSize; x++ {
		fmt.Fprintf(os.Stderr, " %c", labelOfColumns[x+1])
	}
	fmt.Fprintf(os.Stderr, "\n  +")
	for x := 0; x < boardSize; x++ {
		fmt.Fprintf(os.Stderr, "--")
	}
	fmt.Fprintf(os.Stderr, "+\n")
	for y := 0; y < boardSize; y++ {
		fmt.Fprintf(os.Stderr, "%s|", labelOfRows[y+1])
		for x := 0; x < boardSize; x++ {
			fmt.Fprintf(os.Stderr, "%s", stoneLabels[(*board).ColorAtFileRank(x+1, y+1)])
		}
		fmt.Fprintf(os.Stderr, "|")
		if y == 4 {
			fmt.Fprintf(os.Stderr, "  Ko=%s,moves=%d", (*board).GetNameFromTIdx(board.KoIdx), moves)
		}
		fmt.Fprintf(os.Stderr, "\n")
	}
	fmt.Fprintf(os.Stderr, "  +")
	for x := 0; x < boardSize; x++ {
		fmt.Fprintf(os.Stderr, "--")
	}
	fmt.Fprintf(os.Stderr, "+\n")
}

// PrintSgf - SGF形式の棋譜表示。
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
