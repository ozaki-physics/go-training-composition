package money

type Dollar struct {
	// 可視性をプライベートにすることで ファクトリーメソッドの意味が生まれる
	amout int
}

// CreateDollar
// コンストラクタ(ファクトリーメソッド) DDD を意識して生成メソッドを作った
// ドメイン層 モデル のイメージ
func CreateDollar(amout int) Dollar {
	if amout < 0 {
		// 本当は nil と err を返したいが TDD から脱線するので実装しない
		return Dollar{amout: 0}
	}
	return Dollar{amout: amout}
}

// Amout 構造体の値をプライベートにしたから値を取り出すメソッドが必要
func (d *Dollar) Amout() int {
	return d.amout
}

func (d *Dollar) Times(multiplier int) Dollar {
	return CreateDollar(d.amout * multiplier)
}

func (d *Dollar) Equals(object interface{}) bool {
	// 型アサーション
	// See Go言語の型とreflect: https://qiita.com/atsaki/items/3554f5a0609c59a3e10d
	// var dollar Dollar= object.(Dollar)
	dollar, ok := object.(Dollar)
	if ok == true {
		return d.amout == dollar.Amout()
	} else {
		return false
	}
	// 他にも参考になりそうなドキュメント
	// See Package reflect: https://pkg.go.dev/reflect@go1.17.6
}
