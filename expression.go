package main

type Expression interface {
	accept(Visitor) (any, error)
}

type Literal struct {
	Value interface{}
}

func CreateLiteral(value interface{}) *Literal {
	return &Literal{
		Value: value,
	}
}
func (l *Literal) accept(v Visitor) (any, error) {
	return v.VisitLiteral(l), nil
}

type Unary struct {
	Right   Expression
	Operand *Token
}

func CreateUnary(right Expression, operand *Token) *Unary {
	return &Unary{
		Right:   right,
		Operand: operand,
	}
}
func (u *Unary) accept(v Visitor) (any, error) {
	return v.VisitUnary(u)
}

type Binary struct {
	Left     Expression
	Operator Token
	Right    Expression
}

func CreateBinary(left Expression, operator Token, right Expression) *Binary {
	return &Binary{
		Left:     left,
		Operator: operator,
		Right:    right,
	}
}
func (b *Binary) accept(v Visitor) (any, error) {
	return v.VisitBinary(b)
}

type Ternary struct {
	Left   Expression
	Center Expression
	Right  Expression
}

func CreateTernary(left Expression, center Expression, right Expression) *Ternary {
	return &Ternary{
		Left:   left,
		Center: center,
		Right:  right,
	}
}
func (t *Ternary) accept(v Visitor) (any, error) {
	return v.VisitTernary(t)
}

type Operand struct {
	Value string
}

func CreateOperand(value string) *Operand {
	return &Operand{
		Value: value,
	}
}

type Grouping struct {
	Expression Expression
}

func CreateGroup(expr Expression) *Grouping {
	return &Grouping{
		Expression: expr,
	}
}
func (g *Grouping) accept(v Visitor) (any, error) {
	return v.VisitGrouping(g)
}
