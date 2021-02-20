package phase

// Phase - 黒手番、または白手番。
type Phase int

// state
const (
	// None - 開始。
	None Phase = iota
	// Black - 自分のアカウント名を入力しました
	Black
	// White - 自分のパスワードを入力し、そしてプロンプトを待っています
	White
)
