package presenter

import (
	"fmt"
	"strings"

	be "github.com/muzudho/kifuwarabe-go-base/entities"
)

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

// CreateBoardHeader - 手数などを表示
func CreateBoardHeader(position *be.Position, movesNum int) string {
	return fmt.Sprintf("[ Ko=%s MovesNum=%d ]\n", (*position).GetNameFromTIdx(position.KoIdx), movesNum)
}

// CreateBoardString - 盤を描画
func CreateBoardString(position *be.Position) string {
	boardSize := (*position).BoardSize()

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
			b.WriteString(fmt.Sprintf("%s", stoneLabels[(*position).ColorAtFileRank(x+1, y+1)]))
		}
		b.WriteString("|\n")
	}
	b.WriteString("  +")
	for x := 0; x < boardSize; x++ {
		b.WriteString("--")
	}
	b.WriteString("+\n")

	return b.String()
}

// PrintSgf - SGF形式の棋譜表示
func PrintSgf(position *be.Position, movesNum int, record []int) {
	boardSize := position.BoardSize()

	fmt.Printf("(;GM[1]SZ[%d]KM[%.1f]PB[]PW[]\n", boardSize, position.Komi())
	for i := 0; i < movesNum; i++ {
		tIdx := record[i]
		y := tIdx / position.SentinelWidth()
		x := tIdx - y*position.SentinelWidth()
		var sStone = [2]string{"B", "W"}
		fmt.Printf(";%s", sStone[i&1])
		if tIdx == 0 {
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
func GetPointName(position *be.Position, tIdx int) string {
	if tIdx == 0 {
		return "pass"
	}

	y := tIdx / (*position).SentinelWidth()
	x := tIdx - y*(*position).SentinelWidth()

	ax := labelOfColumns[x]

	return fmt.Sprintf("%c%d", ax, y)
}
