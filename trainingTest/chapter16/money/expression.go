package money

type Expression interface {
	Reduce(*Bank, string) Money
	Puls(Expression) Expression
	Times(int) Expression
}
