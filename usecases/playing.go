package usecases

import (
	"fmt"
	"os"
	"time"

	e "github.com/muzudho/kifuwarabe-gtp/entities"
)

// PlayComputerMove - コンピューター・プレイヤーの指し手。 GoGoV9a から呼び出されます。
func PlayComputerMove(board e.IBoard, color int, fUCT int, printBoardType1 func(e.IBoard), printBoardType2 func(e.IBoard, int)) int {
	var tIdx int
	st := time.Now()
	e.AllPlayouts = 0
	tIdx = board.PrimitiveMonteCalro(color, printBoardType1)
	sec := time.Since(st).Seconds()
	fmt.Fprintf(os.Stderr, "%.1f sec, %.0f playout/sec, play=%s,moves=%d,color=%d,playouts=%d,fUCT=%d\n",
		sec, float64(e.AllPlayouts)/sec, board.GetNameFromTIdx(tIdx), e.MovesCount, color, e.AllPlayouts, fUCT)
	board.AddMoves(tIdx, color, sec, printBoardType2)
	return tIdx
}

// UndoV9 - 一手戻します。
func UndoV9() {
	// Unimplemented.
}
