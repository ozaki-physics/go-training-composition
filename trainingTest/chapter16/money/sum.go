package money

type Sum struct {
	Augend Expression
	Addend Expression
}

// Reduce リデュース 縮小するとかのニュアンス
// Sum 型 では 2つの Money が含まれているから 圧縮して1個の Money になるニュアンスは遠からず
func (s Sum) Reduce(b *Bank, to string) Money {
	// Augend, Addend が インタフェースになったおかげで Sum 型 か Money 型 のいずれかが入る可能性があり
	// もし Sum 型 なら Augend.Reduce() で Sum 型 の Reduce() を再帰的に呼び出すことになる
	amount := s.Augend.Reduce(b, to).Amount() + s.Addend.Reduce(b, to).Amount()
	return CreateMoney(amount, to)
}

func (s Sum) Puls(addend Expression) Expression {
	return Sum{s, addend}
}

func (s Sum) Times(multiplier int) Expression {
	return Sum{s.Augend.Times(multiplier), s.Addend.Times(multiplier)}
}
