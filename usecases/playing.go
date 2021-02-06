package usecases

import (
	"fmt"
	"os"
	"time"

	e "github.com/muzudho/kifuwarabe-uec12/entities"
	p "github.com/muzudho/kifuwarabe-uec12/presenter"
)

// SelfplayV9 - コンピューター同士の対局。
func SelfplayV9(board e.IBoard, printBoardType1 func(e.IBoard), printBoardType2 func(e.IBoard, int)) {
	color := 1

	for {
		fUCT := 1
		if color == 1 {
			fUCT = 0
		}
		tIdx := board.GetComputerMove(color, fUCT, printBoardType1)
		board.AddMovesType1(tIdx, color, printBoardType2)
		// パスで２手目以降で棋譜の１つ前（相手）もパスなら終了します。
		if tIdx == 0 && 1 < e.Moves && e.Record[e.Moves-2] == 0 {
			break
		}
		// 自己対局は300手で終了します。
		if 300 < e.Moves {
			break
		} // too long
		color = e.FlipColor(color)
	}

	p.PrintSgf(board, e.Moves, e.Record)
}

// TestPlayoutV9 - 試しにプレイアウトする。
func TestPlayoutV9(board e.IBoard, printBoardType1 func(e.IBoard), printBoardType2 func(e.IBoard, int)) {
	e.FlagTestPlayout = 1
	board.Playout(1, printBoardType1)
	printBoardType2(board, e.Moves)
	p.PrintSgf(board, e.Moves, e.Record)
}

// PlayComputerMoveV9a - コンピューター・プレイヤーの指し手。 GoGoV9a から呼び出されます。
func PlayComputerMoveV9a(board e.IBoard, color int, fUCT int, printBoardType1 func(e.IBoard), printBoardType2 func(e.IBoard, int)) int {
	var tIdx int
	st := time.Now()
	e.AllPlayouts = 0
	tIdx = board.PrimitiveMonteCalro(color, printBoardType1)
	sec := time.Since(st).Seconds()
	fmt.Fprintf(os.Stderr, "%.1f sec, %.0f playout/sec, play_z=%04d,moves=%d,color=%d,playouts=%d,fUCT=%d\n",
		sec, float64(e.AllPlayouts)/sec, board.GetZ4(tIdx), e.Moves, color, e.AllPlayouts, fUCT)
	board.AddMovesType2(tIdx, color, sec, printBoardType2)
	return tIdx
}

// TestPlayoutV9a - 試しにプレイアウトする。
func TestPlayoutV9a(board e.IBoard, printBoardType1 func(e.IBoard), printBoardType2 func(e.IBoard, int)) {
	e.FlagTestPlayout = 1
	board.Playout(1, printBoardType1)
	printBoardType2(board, e.Moves)
	p.PrintSgf(board, e.Moves, e.Record)
}

// UndoV9 - 一手戻します。
func UndoV9() {
	// Unimplemented.
}
