package money

type Sum struct {
	Augend Expression
	Addend Expression
}

// Reduce リデュース 縮小するとかのニュアンス
// Sum 型 では 2つの Money が含まれているから 圧縮して1個の Money になるニュアンスは遠からず
func (s Sum) Reduce(b *Bank, to string) Money {
	amount := s.Augend.Reduce(b, to).Amount() + s.Addend.Reduce(b, to).Amount()
	return CreateMoney(amount, to)
}

func (s Sum) Puls(addend Expression) Expression {
	return Sum{}
}
