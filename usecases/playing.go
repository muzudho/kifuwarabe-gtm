package usecases

import (
	"fmt"
	"os"
	"time"

	be "github.com/muzudho/kifuwarabe-go-base/entities"
	e "github.com/muzudho/kifuwarabe-gtp/entities"
)

// PlayComputerMove - コンピューター・プレイヤーの指し手。 main から呼び出されます。
func PlayComputerMove(position *be.Position, color int, fUCT int, printBoard func(*be.Position)) int {
	var tIdx int
	st := time.Now()
	e.AllPlayouts = 0
	tIdx = e.PrimitiveMonteCalro(position, color, printBoard)
	sec := time.Since(st).Seconds()
	fmt.Fprintf(os.Stderr, "%.1f sec, %.0f playout/sec, play=%s,moves=%d,color=%d,playouts=%d,fUCT=%d\n",
		sec, float64(e.AllPlayouts)/sec, (*position).GetNameFromTIdx(tIdx), position.MovesNum, color, e.AllPlayouts, fUCT)

	// TODO サーバーから返ってきた時刻ではなく、自己計測の時間を入れてる？
	(*position).AddMoves(tIdx, color, sec)

	return tIdx
}

// UndoV9 - 一手戻します。
func UndoV9() {
	// Unimplemented.
}
