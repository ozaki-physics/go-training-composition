package money

type Dollar struct {
	Amout int
}

// CreateDollar コンストラクタ DDD を意識して生成メソッドを作った
func CreateDollar(amout int) Dollar {
	return Dollar{Amout: amout}
}

func (d *Dollar) Times(multiplier int) Dollar {
	return CreateDollar(d.Amout * multiplier)
}
