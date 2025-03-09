package main

type Statement interface {
	accept(StatementVisitor) error
}

type StatementVisitor interface {
	VisitExpressionStatement(e *ExpressionStatement) error
	VisitPrintStatement(p *PrintStatement) error
}

type ExpressionStatement struct {
	Expr Expression
}

func (e *ExpressionStatement) accept(visitor StatementVisitor) error {
	return visitor.VisitExpressionStatement(e)
}

func CreateExpressionStatement(expr Expression) *ExpressionStatement {
	return &ExpressionStatement{
		Expr: expr,
	}
}

type PrintStatement struct {
	Expr Expression
}

func (p *PrintStatement) accept(visitor StatementVisitor) error {
	return visitor.VisitPrintStatement(p)
}

func CreatePrintStatement(expr Expression) *PrintStatement {
	return &PrintStatement{
		Expr: expr,
	}
}
