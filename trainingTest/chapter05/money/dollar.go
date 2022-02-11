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

func (d *Dollar) Times(multiplier int) Dollar {
	return CreateDollar(d.amout * multiplier)
}

// Equals Go 言語では 構造体の中の値で比較してくれるため 実際は実装する必要ない
// 今回は TDD の練習のため 実装する
func (d *Dollar) Equals(object interface{}) bool {
	// 型アサーション
	// See Go言語の型とreflect: https://qiita.com/atsaki/items/3554f5a0609c59a3e10d
	dollar, ok := object.(Dollar)
	if ok == true {
		return d.amout == dollar.amout
	} else {
		return false
	}
	// 他にも参考になりそうなドキュメント
	// See Package reflect: https://pkg.go.dev/reflect@go1.17.6
}
