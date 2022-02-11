package money

// Accessor インタフェースを定義する意味は
// Equals() の引数で 型一致 から インタフェースを実装しているか と制限を緩くするため
type Accessor interface {
	Amount() int
	Currency() string
}

type Money struct {
	// 可視性をプライベートにすることで setter が存在しないことの意味が生まれる
	amount       int
	currencyCode string
}

func (m Money) Amount() int {
	return m.amount
}

func (m Money) Currency() string {
	return m.currencyCode
}

func CreateMoney(a int, c string) Money {
	return Money{a, c}
}

// Equals 引数をインタフェースにすることで
// インタフェースを実装してある構造体なら渡せるようになる
func (m *Money) Equals(a Accessor) bool {
	equal01 := m.amount == a.Amount()
	// 書籍は型を取得して 一致するか比較しているが
	// Accessor インタフェースに 通貨コードを格納して比較する
	equal02 := m.currencyCode == a.Currency()
	return equal01 && equal02
}

// Times Dollar と Franc の struct から移植しつつ統合
func (m *Money) Times(multiplier int) Money {
	return Money{m.Amount() * multiplier, m.Currency()}
}

// CreateDollar Factory Method パターン
func CreateDollar(v int) Money {
	if v < 0 {
		// 本当は DDD 的には nil と err を返したいが TDD から脱線するので実装しない
		return CreateMoney(0, "USD")
	}
	return CreateMoney(v, "USD")
}

// CreateFranc Factory Method パターン
func CreateFranc(v int) Money {
	if v < 0 {
		// 本当は DDD 的には nil と err を返したいが TDD から脱線するので実装しない
		return CreateMoney(0, "CHF")
	}
	return CreateMoney(v, "CHF")
}
