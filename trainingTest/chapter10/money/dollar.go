package money

// Dollar 埋め込みをすることで 継承のように Money に定義されているメソッドを借りて使えるようになる
type Dollar struct {
	Money
}
