package money

import (
	"fmt"
)

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

// Puls 通貨の足し算
func (m Money) Puls(addend Money) Expression {
	return Sum{m, addend}
}

// Reduce リデュース 縮小するとかのニュアンス
// Reduce メソッドは Sum 型 で必要だったから作られて
// Bank 型からみたら 引数が Sum 型 だろうと Money 型だろうと意識しないで済むように
// インタフェースを実装するために Money 型 にも定義された
func (m Money) Reduce(to string) Money {
	return m
}

// String 出力するときのフォーマットを定義する
// レシーバをポインタにすると ポインタの出力だけにフォーマットが適用される
func (m Money) String() string {
	// return fmt.Sprintf("%+v\n", &m)
	// 構造体をそのまま出力しても良かったが それだと 意図するものとフィールドの値が異なっていても気づかない可能性がある
	// 意図的に値を取り出して出力する
	return fmt.Sprintf("{Amount: %v, Currency: %v}", m.Amount(), m.Currency())
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
