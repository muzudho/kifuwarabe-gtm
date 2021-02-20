package entities

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strings"
)

// IBoard - 盤
type IBoard interface {
	// 盤について
	CopyData() []int
	ImportData(boardCopy2 []int)
	BoardSize() int
	// SentinelWidth - 枠付きの盤の一辺の交点数
	SentinelWidth() int
	SentinelBoardMax() int

	// 交点について
	// 指定した交点の石の色
	ColorAt(tIdx int) int
	ColorAtXy(x int, y int) int
	SetColor(tIdx int, color int)
	Exists(tIdx int) bool
	// 石を置きます。
	PutStone(tIdx int, color int, fillEyeErr int) int
	// GetTIdxFromXy - YX形式の座標を、tIdx（配列のインデックス）へ変換します。
	GetTIdxFromXy(x int, y int) int
	// GetNameFromTIdx -
	GetNameFromTIdx(tIdx int) string

	// 検索について
	GetEmptyTIdx() int

	// 呼吸点について
	CountLiberty(tIdx int, pLiberty *int, pStone *int)

	// ゲーム操作について
	// AddMoves - 指し手の追加？
	AddMoves(tIdx int, color int, sec float64, printBoardType2 func(IBoard, int))
	TakeStone(tIdx int, color int)

	// ゲームについて
	// 6.5 といった数字を入れるだけ。実行速度優先で 64bitに。
	Komi() float64
	MaxMoves() int
}

// Record - 棋譜
var Record []int

// RecordTime - 一手にかかった時間。
var RecordTime []float64

// Dir4 - ４方向（右、下、左、上）の番地。初期値は仮の値。
var Dir4 = [4]int{1, 9, -1, 9}

const (
	// DoNotFillEye - 自分の眼を埋めるなってこと☆（＾～＾）
	DoNotFillEye = 1
	// MayFillEye - 自分の眼を埋めてもいいってこと☆（＾～＾）
	MayFillEye = 0
)

// KoIdx - コウの交点。Idx（配列のインデックス）表示。 0 ならコウは無し？
var KoIdx int

// For count liberty.
var checkBoard = []int{}

// MovesCount - 手数
var MovesCount int

var labelOfColumns = []string{"A", "B", "C", "D", "E", "F", "G", "H", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T"}

// Board - 盤。
type Board struct {
	data             []int
	boardSize        int
	sentinelWidth    int
	sentinelBoardMax int
	komi             float64
	maxMoves         int
}

// BoardSize - 何路盤か
func (board Board) BoardSize() int {
	return board.boardSize
}

// SentinelWidth - 枠付きの盤の一辺の交点数
func (board Board) SentinelWidth() int {
	return board.sentinelWidth
}

// SentinelBoardMax - 枠付きの盤の交点数
func (board Board) SentinelBoardMax() int {
	return board.sentinelBoardMax
}

// Komi - コミ
func (board Board) Komi() float64 {
	return board.komi
}

// MaxMoves - 最大手数
func (board Board) MaxMoves() int {
	return board.maxMoves
}

// ColorAt - 指定した交点の石の色
func (board Board) ColorAt(z int) int {
	return board.data[z]
}

// ColorAtXy - 指定した交点の石の色
func (board Board) ColorAtXy(x int, y int) int {
	return board.data[(y+1)*board.sentinelWidth+x+1]
}

// Exists - 指定の交点に石があるか？
func (board Board) Exists(tIdx int) bool {
	return board.data[tIdx] != 0
}

// SetColor - 盤データ。
func (board *Board) SetColor(tIdx int, color int) {
	board.data[tIdx] = color
}

// CopyData - 盤データのコピー。
func (board Board) CopyData() []int {
	boardMax := board.SentinelBoardMax()

	var boardCopy2 = make([]int, boardMax)
	copy(boardCopy2[:], board.data[:])
	return boardCopy2
}

// ImportData - 盤データのコピー。
func (board *Board) ImportData(boardCopy2 []int) {
	copy(board.data[:], boardCopy2[:])
}

/*
// GetZ4 - tIdx（配列のインデックス）を XXYY形式（3～4桁の数）の座標へ変換します。
func (board Board) GetZ4(tIdx int) int {
	if tIdx == 0 {
		return 0
	}
	y := tIdx / board.SentinelWidth()
	x := tIdx - y*board.SentinelWidth()
	return x*100 + y
}
*/

// GetNameFromTIdx -
func (board Board) GetNameFromTIdx(tIdx int) string {
	x, y := board.GetXYFromTIdx(tIdx)
	return GetNameFromXY(x, y)
}

// GetNameFromXY - (1,1) を "A1" に変換
func GetNameFromXY(x int, y int) string {
	return fmt.Sprintf("%s%d", labelOfColumns[x], y)
}

// GetTIdxFromXy - x,y を tIdx（配列のインデックス）へ変換します。
func (board Board) GetTIdxFromXy(x int, y int) int {
	return (y+1)*board.SentinelWidth() + x + 1
}

// GetXYFromTIdx - x,y を tIdx（配列のインデックス）へ変換します。
func (board Board) GetXYFromTIdx(tIdx int) (int, int) {
	return tIdx / board.SentinelWidth(), tIdx % board.SentinelWidth()
}

// GetEmptyTIdx - 空点の tIdx（配列のインデックス）を返します。
func (board Board) GetEmptyTIdx() int {
	var x, y, tIdx int
	for {
		// ランダムに交点を選んで、空点を見つけるまで繰り返します。
		x = rand.Intn(9)
		y = rand.Intn(9)
		tIdx = board.GetTIdxFromXy(x, y)
		if !board.Exists(tIdx) {
			break
		}
	}
	return tIdx
}

// GetXYFromName - "A1" を (1,1) に変換します
func GetXYFromName(name string) (int, int, error) {
	if name == "Pass" {
		return 0, 0, nil
	}

	regexCoord := *regexp.MustCompile("([A-Za-z])(\\d+)")
	matches211 := regexCoord.FindSubmatch([]byte(name))

	var xStr string
	var yStr string
	if 1 < len(matches211) {
		xStr = strings.ToUpper(string(matches211[1]))
		yStr = string(matches211[2])
	} else {
		message := fmt.Sprintf("Unexpected name=[%s]", name)
		return 0, 0, errors.New(message)
	}

	var x int
	switch xStr {
	case "A":
		x = 0
	case "B":
		x = 1
	case "C":
		x = 2
	case "D":
		x = 3
	case "E":
		x = 4
	case "F":
		x = 5
	case "G":
		x = 6
	case "H":
		x = 7
	case "J":
		x = 8
	case "K":
		x = 9
	case "L":
		x = 10
	case "M":
		x = 11
	case "N":
		x = 12
	case "O":
		x = 13
	case "P":
		x = 14
	case "Q":
		x = 15
	case "R":
		x = 16
	case "S":
		x = 17
	case "T":
		x = 18
	default:
		message := fmt.Sprintf("Unexpected xStr=[%s]", xStr)
		return 0, 0, errors.New(message)
	}

	var y int
	switch yStr {
	case "1":
		y = 0
	case "2":
		y = 1
	case "3":
		y = 2
	case "4":
		y = 3
	case "5":
		y = 4
	case "6":
		y = 5
	case "7":
		y = 6
	case "8":
		y = 7
	case "9":
		y = 8
	case "10":
		y = 9
	case "11":
		y = 10
	case "12":
		y = 11
	case "13":
		y = 12
	case "14":
		y = 13
	case "15":
		y = 14
	case "16":
		y = 15
	case "17":
		y = 16
	case "18":
		y = 17
	case "19":
		y = 18
	default:
		message := fmt.Sprintf("Unexpected yStr=[%s]", yStr)
		return 0, 0, errors.New(message)
	}

	return x, y, nil
}

// NewBoard - 盤を作成します。
func NewBoard(data []int, boardSize int, sentinelBoardMax int, komi float64, maxMoves int) *Board {
	board := new(Board)
	board.data = data
	board.boardSize = boardSize
	board.sentinelWidth = boardSize + 2
	board.sentinelBoardMax = sentinelBoardMax
	board.komi = komi
	board.maxMoves = maxMoves

	checkBoard = make([]int, board.SentinelBoardMax())
	Record = make([]int, board.MaxMoves())
	RecordTime = make([]float64, board.MaxMoves())
	Dir4 = [4]int{1, board.SentinelWidth(), -1, -board.SentinelWidth()}

	return board
}

// FlipColor - 白黒反転させます。
func FlipColor(col int) int {
	return 3 - col
}

func (board Board) countLibertySub(tIdx int, color int, pLiberty *int, pStone *int) {
	checkBoard[tIdx] = 1
	*pStone++
	for i := 0; i < 4; i++ {
		z := tIdx + Dir4[i]
		if checkBoard[z] != 0 {
			continue
		}
		if !board.Exists(z) {
			checkBoard[z] = 1
			*pLiberty++
		}
		if board.data[z] == color {
			board.countLibertySub(z, color, pLiberty, pStone)
		}
	}
}

// CountLiberty - 呼吸点を数えます。
func (board Board) CountLiberty(tIdx int, pLiberty *int, pStone *int) {
	*pLiberty = 0
	*pStone = 0
	boardMax := board.SentinelBoardMax()
	// 初期化
	for tIdx2 := 0; tIdx2 < boardMax; tIdx2++ {
		checkBoard[tIdx2] = 0
	}
	board.countLibertySub(tIdx, board.data[tIdx], pLiberty, pStone)
}

// TakeStone - 石を打ち上げ（取り上げ、取り除き）ます。
func (board *Board) TakeStone(tIdx int, color int) {
	board.data[tIdx] = 0
	for dir := 0; dir < 4; dir++ {
		tIdx2 := tIdx + Dir4[dir]
		if board.data[tIdx2] == color {
			board.TakeStone(tIdx2, color)
		}
	}
}

// InitBoard - 盤の初期化
func (board *Board) InitBoard() {
	boardMax := board.SentinelBoardMax()
	boardSize := board.BoardSize()

	// 盤を 枠線　で埋めます
	for tIdx := 0; tIdx < boardMax; tIdx++ {
		board.SetColor(tIdx, 3)
	}

	// 盤上に石を置きます
	for y := 0; y < boardSize; y++ {
		for x := 0; x < boardSize; x++ {
			board.SetColor(board.GetTIdxFromXy(x, y), 0)
		}
	}

	MovesCount = 0
	KoIdx = 0
}

// PutStone - 石を置きます
func (board *Board) PutStone(tIdx int, color int, fillEyeErr int) int {
	var around = [4][3]int{}
	var liberty, stone int
	unCol := FlipColor(color)
	space := 0
	wall := 0
	mycolSafe := 0
	captureSum := 0
	koMaybe := 0

	if tIdx == 0 {
		KoIdx = 0
		return 0
	}
	for dir := 0; dir < 4; dir++ {
		around[dir][0] = 0
		around[dir][1] = 0
		around[dir][2] = 0
		z := tIdx + Dir4[dir]
		color2 := board.ColorAt(z)
		if color2 == 0 {
			space++
		}
		if color2 == 3 {
			wall++
		}
		if color2 == 0 || color2 == 3 {
			continue
		}
		board.CountLiberty(z, &liberty, &stone)
		around[dir][0] = liberty
		around[dir][1] = stone
		around[dir][2] = color2
		if color2 == unCol && liberty == 1 {
			captureSum += stone
			koMaybe = z
		}
		if color2 == color && 2 <= liberty {
			mycolSafe++
		}

	}
	if captureSum == 0 && space == 0 && mycolSafe == 0 {
		return 1
	}
	if tIdx == KoIdx {
		return 2
	}
	if wall+mycolSafe == 4 && fillEyeErr == DoNotFillEye {
		return 3
	}
	if board.Exists(tIdx) {
		return 4
	}

	for dir := 0; dir < 4; dir++ {
		lib := around[dir][0]
		color2 := around[dir][2]
		if color2 == unCol && lib == 1 && board.Exists(tIdx+Dir4[dir]) {
			board.TakeStone(tIdx+Dir4[dir], unCol)
		}
	}

	board.SetColor(tIdx, color)

	board.CountLiberty(tIdx, &liberty, &stone)

	if captureSum == 1 && stone == 1 && liberty == 1 {
		KoIdx = koMaybe
	} else {
		KoIdx = 0
	}
	return 0
}

// AddMoves - 指し手の追加？
func (board *Board) AddMoves(tIdx int, color int, sec float64, printBoardType2 func(IBoard, int)) {
	err := board.PutStone(tIdx, color, MayFillEye)
	if err != 0 {
		fmt.Fprintf(os.Stderr, "(AddMoves) Err=%d\n", err)
		os.Exit(0)
	}
	Record[MovesCount] = tIdx
	RecordTime[MovesCount] = sec
	MovesCount++
	printBoardType2(board, MovesCount)
}
