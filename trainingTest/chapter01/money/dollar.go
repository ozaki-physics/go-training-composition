package money

type Dollar struct {
	Amout int
}

func (d *Dollar) Times(multiplier int) {
	d.Amout *= multiplier
}
