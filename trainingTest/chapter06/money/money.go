package money

// AmountGetter インタフェースを定義する意味は
// Equals() の引数で 型一致 から インタフェースを実装しているか と制限を緩くするため
// このインタフェースは private にしても動作するが インタフェースという性質上 public にしておく
type AmountGetter interface {
	// 実装を強制するメソッドは private でもよい
	getAmount() int
}

type Money struct {
	// 可視性をプライベートにすることで ファクトリーメソッドの意味が生まれる
	amount int
}

// getAmount
// レシーバをポインタにすると テストが通らなくなる
// なぜなら ポインタにすると 埋め込みの子(money.Dollar)が インタフェース(money.AmountGetter)を
// 実装してないことになるから 埋め込みの子は Equals() の引数になれないと言われる
// なんか納得できないけど そういうものと認識して一旦置いとこう
func (m Money) getAmount() int {
	return m.amount
}

// Equals 引数をインタフェースにすることで
// getAmount() メソッドを持っている構造体なら許されるようになる
// つまり Money 型しか受け付けない -> Dollar, Franc(埋め込みで getAmount() を持っている)が許される
func (m *Money) Equals(a AmountGetter) bool {
	return m.getAmount() == a.getAmount()
}
