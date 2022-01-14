package money

type Franc struct {
	amout int
}

func CreateFranc(amout int) Franc {
	if amout < 0 {
		// 本当は nil と err を返したいが TDD から脱線するので実装しない
		return Franc{amout: 0}
	}
	return Franc{amout: amout}
}

func (d *Franc) Times(multiplier int) Franc {
	return CreateFranc(d.amout * multiplier)
}

func (d *Franc) Equals(object interface{}) bool {
	// 型アサーション
	franc, ok := object.(Franc)
	if ok == true {
		return d.amout == franc.amout
	} else {
		return false
	}
}
