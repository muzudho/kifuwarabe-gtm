package stone

// Stone - 石。
type Stone int

// state
const (
	// None - 開始。
	None Stone = iota
	// Black - 自分のアカウント名を入力しました
	Black
	// White - 自分のパスワードを入力し、そしてプロンプトを待っています
	White
	// Wall - 壁
	Wall
)

// FlipColor - 白黒反転させます。
func FlipColor(col int) int {
	return 3 - col
}
