package money

// Accessor インタフェースを定義する意味は
// Equals() の引数で 型一致 から インタフェースを実装しているか と制限を緩くするため
type Accessor interface {
	amount() int
	name() string // Chapter が進むと currencyCode になるが 書籍の意図を組んで まだ currencyCode にしない
}

type Money struct {
	// 可視性をプライベートにすることで setter が存在しないことの意味が生まれる
	volume       int
	currencyCode string
}

func (m Money) amount() int {
	return m.volume
}

func (m Money) name() string {
	return m.currencyCode
}

// Equals 引数をインタフェースにすることで
// インタフェースを実装してある構造体なら渡せるようになる
func (m *Money) Equals(a Accessor) bool {
	equal01 := m.amount() == a.amount()
	// 書籍は型を取得して 一致するか比較しているが
	// Accessor インタフェースに 通貨コードを格納して比較する
	equal02 := m.name() == a.name()
	return equal01 && equal02
}

// Times Dollar と Franc の struct から移植しつつ統合
func (m *Money) Times(multiplier int) Money {
	return Money{volume: m.amount() * multiplier, currencyCode: m.name()}
}

// CreateDollar Factory Method パターン
func CreateDollar(v int) Money {
	if v < 0 {
		// 本当は DDD 的には nil と err を返したいが TDD から脱線するので実装しない
		return Money{volume: 0, currencyCode: "Dollar"}
	}
	return Money{volume: v, currencyCode: "Dollar"}
}

// CreateFranc Factory Method パターン
func CreateFranc(v int) Money {
	if v < 0 {
		// 本当は DDD 的には nil と err を返したいが TDD から脱線するので実装しない
		return Money{volume: 0, currencyCode: "Franc"}
	}
	return Money{volume: v, currencyCode: "Franc"}
}
