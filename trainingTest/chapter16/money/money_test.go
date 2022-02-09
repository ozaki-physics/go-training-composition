package money_test

import (
	"fmt"
	"testing"

	. "github.com/ozaki-physics/go-training-composition/trainingTest/chapter16/money"
)

func TestCreateDollar(t *testing.T) {
	t.Run("コンストラクタ(ファクトリーメソッド)が使える", func(t *testing.T) {
		// Go 言語の仕様として Equals メソッドを作らなくても 値の比較をしてくれる
		if CreateDollar(5) != CreateDollar(5) {
			t.Errorf("ファクトリーメソッドが使えてない or Go 言語の構造体比較の仕様を誤解している")
		}
	})
}

func TestTimes(t *testing.T) {
	t.Run("もとの $5 は変わらず 何度でもかけ算ができる", func(t *testing.T) {
		// 型を明示しても書いてみる
		var five Money = CreateDollar(5)

		ten := CreateDollar(10)
		expected := true
		// Times の戻り値を Expression にしたから Accessor インタフェースを実装していないことになってしまった
		// Expression の実体は Money だから存在はするが...
		// 方法は2つ
		// 1. Expression を型アサーションして Money にする
		// 2. Expression に Accessor インタフェースも加える
		// もし 2つ目を選択したら Sum 型 にも Amount() int, Currency() string の実装が必要になるが
		// Sum 型 は Money が2つあるから どっちの amount か分からないし 合算後の amount にしようとしても
		// Amount() の引数に Bank が渡せないし 難しそう
		// よって 型アサーションを選択する
		// クラスの明示的なチェックは ポリモフィズムに置き換えるべきだ という筆者の声が耳が痛い
		actual := ten.Equals(five.Times(2).(Money))
		if expected != actual {
			t.Errorf("期待値 %v で実際は %v", expected, actual)
		}

		fifteen := CreateDollar(15)
		expected = true
		actual = fifteen.Equals(five.Times(3).(Money))
		if expected != actual {
			t.Errorf("期待値 %v で実際は %v", expected, actual)
		}
	})
}

func TestEquals(t *testing.T) {
	t.Run("同じ金額が等価の定義", func(t *testing.T) {
		var five Money = CreateDollar(5)

		expected := true
		actual := five.Equals(CreateDollar(5))
		if expected != actual {
			t.Errorf("Dollar: 期待値 %v で実際は %v", expected, actual)
		}

		expected = false
		actual = five.Equals(CreateDollar(6))
		if expected != actual {
			t.Errorf("Dollar: 期待値 %v で実際は %v", expected, actual)
		}

		var fiveFranc Money = CreateFranc(5)

		expected = true
		actual = fiveFranc.Equals(CreateFranc(5))
		if expected != actual {
			t.Errorf("Franc: 期待値 %v で実際は %v", expected, actual)
		}

		expected = false
		actual = fiveFranc.Equals(CreateFranc(6))
		if expected != actual {
			t.Errorf("Franc: 期待値 %v で実際は %v", expected, actual)
		}
	})
}

func TestCreateFranc(t *testing.T) {
	t.Run("コンストラクタ(ファクトリーメソッド)が使える", func(t *testing.T) {
		// Go 言語の仕様として Equals メソッドを作らなくても 値の比較をしてくれる
		if CreateFranc(5) != CreateFranc(5) {
			t.Errorf("ファクトリーメソッドが使えてない or Go 言語の構造体比較の仕様を誤解している")
		}
	})
}

func TestFrancTimes(t *testing.T) {
	t.Run("もとの 5CHF は変わらず 何度でもかけ算ができる", func(t *testing.T) {
		// 型を明示しても書いてみる
		var five Money = CreateFranc(5)

		ten := CreateFranc(10)
		expected := true
		actual := ten.Equals(five.Times(2).(Money))
		if expected != actual {
			t.Errorf("Franc: 期待値 %v で実際は %v", expected, actual)
		}

		fifteen := CreateFranc(15)
		expected = true
		actual = fifteen.Equals(five.Times(3).(Money))
		if expected != actual {
			t.Errorf("Franc: 期待値 %v で実際は %v", expected, actual)
		}
	})
}

func TestEqualsDollarFranc(t *testing.T) {
	t.Run("同じ数値の Dollar と Franc を比較して 等価", func(t *testing.T) {
		var five Money = CreateFranc(5)
		expected := false
		actual := five.Equals(CreateDollar(5))
		if expected != actual {
			t.Errorf("Franc: 期待値 %v で実際は %v", expected, actual)
		}
	})
}

func TestCurrency(t *testing.T) {
	t.Run("Dollar と Franc を統合する 通貨(Currency)概念を導入する", func(t *testing.T) {
		var oneUSD Money = CreateDollar(1)
		expected := "USD"
		actual := oneUSD.Currency()
		if expected != actual {
			t.Errorf("Franc: 期待値 %v で実際は %v", expected, actual)
		}
		var oneCHF Money = CreateFranc(1)
		expected = "CHF"
		actual = oneCHF.Currency()
		if expected != actual {
			t.Errorf("Franc: 期待値 %v で実際は %v", expected, actual)
		}
	})
}

func TestString(t *testing.T) {
	t.Run("toString 的なメソッドが動作するか確認する", func(t *testing.T) {
		var oneUSD Money = CreateDollar(1)
		fmt.Println(oneUSD)
		fmt.Println(&oneUSD)
		var oneCHF Money = CreateFranc(1)
		fmt.Println(oneCHF)
		fmt.Println(&oneCHF)
	})
}

func TestSimpleAdditon(t *testing.T) {
	t.Run("ドル同士の足し算ができる", func(t *testing.T) {
		five := CreateDollar(5)
		var sum Expression = five.Puls(five)

		// reduced は Expression に為替レートを適用することで得られる換算結果
		// 現実で 為替レートを使った換算を司るもの = 銀行
		bank := Bank{}
		var reduced Money = bank.Reduce(sum, "USD")
		expected := true
		actual := reduced.Equals(CreateDollar(10))
		if expected != actual {
			t.Errorf("Franc: 期待値 %v で実際は %v", expected, actual)
		}
	})
}

func TestPlusReturnsSum(t *testing.T) {
	t.Run("Plus メソッドの戻り値が Sum か", func(t *testing.T) {
		var five Money = CreateDollar(5)
		var result Expression = five.Puls(five)
		// キャストじゃなくて型アサーション
		var sum Sum = result.(Sum)

		expected := true
		actual := five.Equals(sum.Augend.(Money))
		if expected != actual {
			t.Errorf("Franc: 期待値 %v で実際は %v", expected, actual)
		}

		expected = true
		actual = five.Equals(sum.Addend.(Money))
		if expected != actual {
			t.Errorf("Franc: 期待値 %v で実際は %v", expected, actual)
		}
	})
}

func TestReduceSum(t *testing.T) {
	t.Run("Sum で足される通貨が同じなら 足し算の結果は同じ", func(t *testing.T) {
		var sum Sum = Sum{CreateDollar(3), CreateDollar(4)}

		bank := Bank{}
		var reduced Money = bank.Reduce(sum, "USD")
		expected := true
		actual := reduced.Equals(CreateDollar(7))
		if expected != actual {
			t.Errorf("Franc: 期待値 %v で実際は %v", expected, actual)
		}
	})
}

func TestReduceMoney(t *testing.T) {
	t.Run("bank.Reduce の引数を money にしても reduce に渡す通貨が同じなら 同じ値", func(t *testing.T) {
		bank := Bank{}
		var result Money = bank.Reduce(CreateDollar(1), "USD")

		expected := true
		actual := result.Equals(CreateDollar(1))
		if expected != actual {
			t.Errorf("Franc: 期待値 %v で実際は %v", expected, actual)
		}
	})
}

func TestReduceMoneyDifferentCurrency(t *testing.T) {
	t.Run("2 CHF を 1 USD にする", func(t *testing.T) {
		bank := CreateBank()
		bank.AddRate("CHF", "USD", 2)
		var result Money = bank.Reduce(CreateFranc(2), "USD")

		expected := true
		actual := result.Equals(CreateDollar(1))
		if expected != actual {
			t.Errorf("Franc: 期待値 %v で実際は %v", expected, actual)
		}
	})
}

func TestIdentityRate(t *testing.T) {
	t.Run("USD から USD への為替レートは1", func(t *testing.T) {
		bank := CreateBank()

		expected := 1
		actual := bank.Rate("USD", "USD")
		if expected != actual {
			t.Errorf("Franc: 期待値 %v で実際は %v", expected, actual)
		}
	})
}

func TestMixedAdditon(t *testing.T) {
	t.Run("$5 + 10 CHF", func(t *testing.T) {
		// ネットスラングとして Dollar は Bucks とも言うらしい
		var fiveBucks Expression = CreateDollar(5)
		var tenFrancs Expression = CreateFranc(10)
		bank := CreateBank()
		bank.AddRate("CHF", "USD", 2)
		var result Money = bank.Reduce(fiveBucks.Puls(tenFrancs), "USD")

		expected := true
		actual := result.Equals(CreateDollar(10))
		if expected != actual {
			t.Errorf("Franc: 期待値 %v で実際は %v", expected, actual)
		}
	})
}

func TestSumPlusMoney(t *testing.T) {
	t.Run("($5 + 10 CHF) + $5 = $15 を Sum struct を使って行う", func(t *testing.T) {
		var fiveBucks Expression = CreateDollar(5)
		var tenFrancs Expression = CreateFranc(10)
		bank := CreateBank()
		bank.AddRate("CHF", "USD", 2)

		// テストだから直接的に Sum 型 の struct を生成する
		var sum Expression = Sum{fiveBucks, tenFrancs}.Puls(fiveBucks)
		var result Money = bank.Reduce(sum, "USD")

		expected := true
		actual := result.Equals(CreateDollar(15))
		if expected != actual {
			t.Errorf("Franc: 期待値 %v で実際は %v", expected, actual)
		}
	})
}

func TestSumTimes(t *testing.T) {
	t.Run("($5 + 10 CHF) * 2 = $20 を Sum struct を使って行う", func(t *testing.T) {
		var fiveBucks Expression = CreateDollar(5)
		var tenFrancs Expression = CreateFranc(10)
		bank := CreateBank()
		bank.AddRate("CHF", "USD", 2)

		// テストだから直接的に Sum 型 の struct を生成する
		var sum Expression = Sum{fiveBucks, tenFrancs}.Times(2)
		var result Money = bank.Reduce(sum, "USD")

		expected := true
		actual := result.Equals(CreateDollar(20))
		if expected != actual {
			t.Errorf("Franc: 期待値 %v で実際は %v", expected, actual)
		}
	})
}
