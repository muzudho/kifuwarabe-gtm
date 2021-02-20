package entities

import (
	"fmt"
	"math/rand"
	"time"
)

// AllPlayouts - プレイアウトした回数
var AllPlayouts int

// UctChildrenSize - UCTの最大手数
var UctChildrenSize int

// Playout - 最後まで石を打ちます。得点を返します。
func (board *Board) Playout(turnColor int, printBoard func(*Board)) int {
	boardSize := (*board).BoardSize()

	color := turnColor
	previousTIdx := 0
	loopMax := boardSize*boardSize + 200
	boardMax := (*board).SentinelBoardMax()

	AllPlayouts++
	for loop := 0; loop < loopMax; loop++ {
		var empty = make([]int, boardMax)
		var emptyNum, r, tIdx int
		for y := 0; y <= boardSize; y++ {
			for x := 0; x < boardSize; x++ {
				tIdx = (*board).GetTIdxFromFileRank(x+1, y+1)
				if (*board).Exists(tIdx) {
					continue
				}
				empty[emptyNum] = tIdx
				emptyNum++
			}
		}
		r = 0
		for {
			if emptyNum == 0 {
				tIdx = 0
			} else {
				r = rand.Intn(emptyNum)
				tIdx = empty[r]
			}
			err := (*board).PutStone(tIdx, color, DoNotFillEye)
			if err == 0 {
				break
			}
			empty[r] = empty[emptyNum-1]
			emptyNum--
		}
		if tIdx == 0 && previousTIdx == 0 {
			break
		}
		previousTIdx = tIdx
		// printBoard()
		// fmt.Printf("loop=%d,tIdx=%s,c=%d,emptyNum=%d,Ko=%s\n",
		// 	loop, e.GetNameFromXY(tIdx), color, emptyNum, e.GetNameFromXY(board.KoIdx()))
		color = FlipColor(color)
	}
	return board.countScore(turnColor)
}

func (board *Board) countScore(turnColor int) int {
	var mk = [4]int{}
	var kind = [3]int{0, 0, 0}
	var score, blackArea, whiteArea, blackSum, whiteSum int
	boardSize := (*board).BoardSize()

	for y := 0; y < boardSize; y++ {
		for x := 0; x < boardSize; x++ {
			tIdx := (*board).GetTIdxFromFileRank(x+1, y+1)
			color2 := (*board).ColorAt(tIdx)
			kind[color2]++
			if color2 != 0 {
				continue
			}
			mk[1] = 0
			mk[2] = 0
			for dir := 0; dir < 4; dir++ {
				mk[(*board).ColorAt(tIdx+Dir4[dir])]++
			}
			if mk[1] != 0 && mk[2] == 0 {
				blackArea++
			}
			if mk[2] != 0 && mk[1] == 0 {
				whiteArea++
			}
		}
	}
	blackSum = kind[1] + blackArea
	whiteSum = kind[2] + whiteArea
	score = blackSum - whiteSum
	win := 0
	if 0 < float64(score)-(*board).Komi() {
		win = 1
	}
	if turnColor == 2 {
		win = -win
	} // gogo07

	// fmt.Printf("blackSum=%2d, (stones=%2d, area=%2d)\n", blackSum, kind[1], blackArea)
	// fmt.Printf("whiteSum=%2d, (stones=%2d, area=%2d)\n", whiteSum, kind[2], whiteArea)
	// fmt.Printf("score=%d, win=%d\n", score, win)
	return win
}

// PrimitiveMonteCalro - モンテカルロ木探索 Version 9a.
func (board *Board) PrimitiveMonteCalro(color int, printBoard func(*Board)) int {
	boardSize := (*board).BoardSize()

	// ９路盤なら
	// tryNum := 30
	// １９路盤なら
	tryNum := 3
	bestTIdx := 0
	var bestValue, winRate float64
	var boardCopy = (*board).CopyData()
	koZCopy := board.KoIdx
	bestValue = -100.0

	for y := 0; y <= boardSize; y++ {
		for x := 0; x < boardSize; x++ {
			tIdx := (*board).GetTIdxFromFileRank(x+1, y+1)
			if (*board).Exists(tIdx) {
				continue
			}
			err := (*board).PutStone(tIdx, color, DoNotFillEye)
			if err != 0 {
				continue
			}

			winSum := 0
			for i := 0; i < tryNum; i++ {
				var boardCopy2 = (*board).CopyData()
				koZCopy2 := board.KoIdx

				win := -board.Playout(FlipColor(color), printBoard)

				winSum += win
				board.KoIdx = koZCopy2
				(*board).ImportData(boardCopy2)
			}
			winRate = float64(winSum) / float64(tryNum)
			if bestValue < winRate {
				bestValue = winRate
				bestTIdx = tIdx
				// fmt.Printf("(primitiveMonteCalroV9) bestTIdx=%s,color=%d,v=%5.3f,tryNum=%d\n", bestTIdx, color, bestValue, tryNum)
			}
			board.KoIdx = koZCopy
			(*board).ImportData(boardCopy)
		}
	}
	return bestTIdx
}

// GetComputerMove - コンピューターの指し手。
func GetComputerMove(board *Board, color int, fUCT int, printBoard func(*Board)) int {
	var tIdx int
	start := time.Now()
	AllPlayouts = 0
	tIdx = board.PrimitiveMonteCalro(color, printBoard)
	sec := time.Since(start).Seconds()
	fmt.Printf("(GetComputerMove) %.1f sec, %.0f playout/sec, play=%s,moves=%d,color=%d,playouts=%d,fUCT=%d\n",
		sec, float64(AllPlayouts)/sec, (*board).GetNameFromTIdx(tIdx), MovesNum, color, AllPlayouts, fUCT)
	return tIdx
}
