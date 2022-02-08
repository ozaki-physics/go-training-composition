package money

type Bank struct {
	// 2つの通貨の組をレートに紐付ける
	rates map[pair]int
}

// 通貨の組
type pair struct {
	from, to string
}

// CreateBank コンストラクタ的生成メソッド
func CreateBank() Bank {
	b := Bank{}
	b.rates = make(map[pair]int)
	return b
}

func (b *Bank) Reduce(source Expression, to string) Money {
	return source.Reduce(b, to)
}

// AddRate 書籍でいうと HashMap に put している
func (b *Bank) AddRate(from, to string, rate int) {
	b.rates[pair{from, to}] = rate
}

// Rate テストで同じ通貨なら為替レートは1をテストするために public メソッドになっている
func (b *Bank) Rate(from string, to string) int {
	// from と to が同じ通貨なら 為替レートは1
	if from == to {
		return 1
	}

	// 通貨の組から為替レートを取り出す
	p := pair{from, to}
	return b.rates[p]
}
