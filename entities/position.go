package entities

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strings"

	"github.com/muzudho/kifuwarabe-gtp/entities/stone"
)

const (
	// DoNotFillEye - 自分の眼を埋めるなってこと☆（＾～＾）
	DoNotFillEye = 1
	// MayFillEye - 自分の眼を埋めてもいいってこと☆（＾～＾）
	MayFillEye = 0
)

var labelOfColumns = []string{"0", "A", "B", "C", "D", "E", "F", "G", "H", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T"}

// Position - 局面
// 盤面や、棋譜を含む
type Position struct {
	boardData        []int
	boardSize        int
	sentinelWidth    int
	sentinelBoardMax int
	// For count liberty.
	checkBoard []int
	// KoIdx - コウの交点。Idx（配列のインデックス）表示。 0 ならコウは無し？
	KoIdx int
	// Dir4 - ４方向（右、下、左、上）の番地。初期値は仮の値。
	Dir4 [4]int

	komi     float64
	maxMoves int

	// MovesNum - 手数
	MovesNum int
	// Record - 棋譜
	Record []int

	// RecordTime - 一手にかかった時間。
	RecordTime []float64
}

// NewPosition - 盤を作成します。
func NewPosition(boardData []int, boardSize int, sentinelBoardMax int, komi float64, maxMoves int) *Position {
	position := new(Position)
	position.boardData = boardData
	position.boardSize = boardSize
	position.sentinelWidth = boardSize + 2
	position.sentinelBoardMax = sentinelBoardMax
	position.komi = komi
	position.maxMoves = maxMoves

	boardMax := position.SentinelBoardMax()
	position.checkBoard = make([]int, boardMax)
	position.Record = make([]int, position.MaxMoves())
	position.RecordTime = make([]float64, position.MaxMoves())
	position.Dir4 = [4]int{1, position.SentinelWidth(), -1, -position.SentinelWidth()}

	// 盤を 枠線　で埋めます
	for tIdx := 0; tIdx < boardMax; tIdx++ {
		position.SetColor(tIdx, 3)
	}

	// 盤上に石を置きます
	for y := 0; y < boardSize; y++ {
		for x := 0; x < boardSize; x++ {
			position.SetColor(position.GetTIdxFromFileRank(x+1, y+1), 0)
		}
	}

	position.MovesNum = 0
	position.KoIdx = 0

	return position
}

// ⚡　盤について

// CopyData - 盤データのコピー。
func (position Position) CopyData() []int {
	boardMax := position.SentinelBoardMax()

	var boardCopy2 = make([]int, boardMax)
	copy(boardCopy2[:], position.boardData[:])
	return boardCopy2
}

// ImportData - 盤データのコピー。
func (position *Position) ImportData(boardCopy2 []int) {
	copy(position.boardData[:], boardCopy2[:])
}

// BoardSize - 何路盤か
func (position Position) BoardSize() int {
	return position.boardSize
}

// SentinelWidth - 枠付きの盤の一辺の交点数
func (position Position) SentinelWidth() int {
	return position.sentinelWidth
}

// SentinelBoardMax - 枠付きの盤の交点数
func (position Position) SentinelBoardMax() int {
	return position.sentinelBoardMax
}

// ⚡ 交点について

// ColorAt - 指定した交点の石の色
func (position Position) ColorAt(tIdx int) int {
	return position.boardData[tIdx]
}

// ColorAtFileRank - 指定した交点の石の色
// * `file` - 1 Origin.
// * `rank` - 1 Origin.
func (position Position) ColorAtFileRank(file int, rank int) int {
	return position.boardData[rank*position.sentinelWidth+file]
}

// SetColor - 盤データ
func (position *Position) SetColor(tIdx int, color int) {
	position.boardData[tIdx] = color
}

// Exists - 指定の交点に石があるか？
func (position Position) Exists(tIdx int) bool {
	return position.boardData[tIdx] != 0
}

// PutStone - 石を置きます
func (position *Position) PutStone(tIdx int, color int, fillEyeErr int) int {
	var around = [4][3]int{}
	var liberty, stoneCount int
	unCol := stone.FlipColor(color)
	space := 0
	wall := 0
	mycolSafe := 0
	captureSum := 0
	koMaybe := 0

	if tIdx == 0 {
		position.KoIdx = 0
		return 0
	}
	for dir := 0; dir < 4; dir++ {
		around[dir][0] = 0
		around[dir][1] = 0
		around[dir][2] = 0
		tIdx2 := tIdx + position.Dir4[dir]
		color2 := position.ColorAt(tIdx2)
		if color2 == 0 {
			space++
		}
		if color2 == 3 {
			wall++
		}
		if color2 == 0 || color2 == 3 {
			continue
		}
		position.CountLiberty(tIdx2, &liberty, &stoneCount)
		around[dir][0] = liberty
		around[dir][1] = stoneCount
		around[dir][2] = color2
		if color2 == unCol && liberty == 1 {
			captureSum += stoneCount
			koMaybe = tIdx2
		}
		if color2 == color && 2 <= liberty {
			mycolSafe++
		}

	}
	if captureSum == 0 && space == 0 && mycolSafe == 0 {
		return 1
	}
	if tIdx == position.KoIdx {
		return 2
	}
	if wall+mycolSafe == 4 && fillEyeErr == DoNotFillEye {
		return 3
	}
	if position.Exists(tIdx) {
		return 4
	}

	for dir := 0; dir < 4; dir++ {
		lib := around[dir][0]
		color2 := around[dir][2]
		if color2 == unCol && lib == 1 && position.Exists(tIdx+position.Dir4[dir]) {
			position.TakeStone(tIdx+position.Dir4[dir], unCol)
		}
	}

	position.SetColor(tIdx, color)

	position.CountLiberty(tIdx, &liberty, &stoneCount)

	if captureSum == 1 && stoneCount == 1 && liberty == 1 {
		position.KoIdx = koMaybe
	} else {
		position.KoIdx = 0
	}
	return 0
}

// GetTIdxFromFileRank - x,y を tIdx（配列のインデックス）へ変換します。
func (position Position) GetTIdxFromFileRank(file int, rank int) int {
	return rank*position.SentinelWidth() + file
}

// GetNameFromTIdx -
func (position Position) GetNameFromTIdx(tIdx int) string {
	file, rank := position.GetFileRankFromTIdx(tIdx)
	return GetNameFromFileRank(file, rank)
}

// GetEmptyTIdx - 空点の tIdx（配列のインデックス）を返します。
func (position Position) GetEmptyTIdx() int {
	var x, y, tIdx int
	for {
		// ランダムに交点を選んで、空点を見つけるまで繰り返します。
		x = rand.Intn(9)
		y = rand.Intn(9)
		tIdx = position.GetTIdxFromFileRank(x+1, y+1)
		if !position.Exists(tIdx) {
			break
		}
	}
	return tIdx
}

// CountLiberty - 呼吸点を数えます。
func (position Position) CountLiberty(tIdx int, pLiberty *int, stoneCount *int) {
	*pLiberty = 0
	*stoneCount = 0
	boardMax := position.SentinelBoardMax()
	// 初期化
	for tIdx2 := 0; tIdx2 < boardMax; tIdx2++ {
		position.checkBoard[tIdx2] = 0
	}
	position.countLibertySub(tIdx, position.boardData[tIdx], pLiberty, stoneCount)
}

func (position Position) countLibertySub(tIdx int, color int, pLiberty *int, stoneCount *int) {
	position.checkBoard[tIdx] = 1
	*stoneCount++
	for i := 0; i < 4; i++ {
		tIdx2 := tIdx + position.Dir4[i]
		if position.checkBoard[tIdx2] != 0 {
			continue
		}
		if !position.Exists(tIdx2) {
			position.checkBoard[tIdx2] = 1
			*pLiberty++
		}
		if position.boardData[tIdx2] == color {
			position.countLibertySub(tIdx2, color, pLiberty, stoneCount)
		}
	}
}

// AddMoves - 指し手の追加？
func (position *Position) AddMoves(tIdx int, color int, sec float64) {
	err := (*position).PutStone(tIdx, color, MayFillEye)
	if err != 0 {
		fmt.Fprintf(os.Stderr, "(AddMoves) Err=%d\n", err)
		os.Exit(0)
	}
	position.Record[position.MovesNum] = tIdx
	position.RecordTime[position.MovesNum] = sec
	position.MovesNum++
}

// Komi - コミ
func (position Position) Komi() float64 {
	return position.komi
}

// MaxMoves - 最大手数
func (position Position) MaxMoves() int {
	return position.maxMoves
}

// GetNameFromFileRank - (1,1) を "A1" に変換
func GetNameFromFileRank(file int, rank int) string {
	return fmt.Sprintf("%s%d", labelOfColumns[file], rank)
}

// GetFileRankFromTIdx - tIdx（配列のインデックス）を、file, rank へ変換します。
func (position Position) GetFileRankFromTIdx(tIdx int) (int, int) {
	return tIdx % position.SentinelWidth(), tIdx / position.SentinelWidth()
}

// GetXYFromName - "A1" を (1,1) に変換します
func GetXYFromName(name string) (int, int, error) {
	if name == "pass" {
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

// TakeStone - 石を打ち上げ（取り上げ、取り除き）ます。
func (position *Position) TakeStone(tIdx int, color int) {
	position.boardData[tIdx] = 0
	for dir := 0; dir < 4; dir++ {
		tIdx2 := tIdx + position.Dir4[dir]
		if position.boardData[tIdx2] == color {
			position.TakeStone(tIdx2, color)
		}
	}
}
