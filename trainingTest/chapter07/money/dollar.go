package money

// Dollar 埋め込みをすることで 継承のように Money に定義されているメソッドを借りて使えるようになる
type Dollar struct {
	Money
}

// CreateDollar
// コンストラクタ(ファクトリーメソッド) DDD を意識して生成メソッドを作った
// ドメイン層 モデル のイメージ
func CreateDollar(amount int) Dollar {
	if amount < 0 {
		// 本当は nil と err を返したいが TDD から脱線するので実装しない
		return Dollar{Money{volume: 0, currencyCode: "Dollar"}}
	}
	return Dollar{Money{volume: amount, currencyCode: "Dollar"}}
}

func (d *Dollar) Times(multiplier int) Dollar {
	return CreateDollar(d.amount() * multiplier)
}
