package money

type Franc struct {
	Money
}

func CreateFranc(amount int) Franc {
	if amount < 0 {
		// 本当は nil と err を返したいが TDD から脱線するので実装しない
		return Franc{Money{amount: 0}}
	}
	return Franc{Money{amount: amount}}
}

func (f *Franc) Times(multiplier int) Franc {
	return CreateFranc(f.amount * multiplier)
}
