package money

type Sum struct {
	Augend Money
	Addend Money
}

// Reduce リデュース 縮小するとかのニュアンス
// Sum 型 では 2つの Money が含まれているから 圧縮して1個の Money になるニュアンスは遠からず
func (s Sum) Reduce(to string) Money {
	amount := s.Augend.Amount() + s.Addend.Amount()
	return CreateMoney(amount, to)
}
