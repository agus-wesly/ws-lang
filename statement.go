package main

type Statement interface {
	accept(StatementVisitor) error
}

type StatementVisitor interface {
	VisitExpressionStatement(e *ExpressionStatement) error
	VisitPrintStatement(p *PrintStatement) error
    VisitVarDeclaration(v *VarDeclaration) error
    VisitBlockStatement(v *BlockStatement) error
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

type VarDeclaration struct {
    Identifier Token
	Expr Expression
}

func (v *VarDeclaration) accept(visitor StatementVisitor) error {
	return visitor.VisitVarDeclaration(v)
}

func CreateVarDeclaration(expr Expression, identifier Token) *VarDeclaration {
	return &VarDeclaration{
		Expr: expr,
        Identifier: identifier,
	}
}

type VarAssignment struct {
    Token Token
	Expr Expression
}

func (v *VarAssignment) accept(visitor ExpressionVisitor) (any, error) {
	return visitor.VisitVarAssignment(v)
}

func CreateVarAssignment(token Token, expr Expression) *VarAssignment {
	return &VarAssignment{
        Token: token,
        Expr: expr,
	}
}

type BlockStatement struct {
    Statements []Statement
}

func (b *BlockStatement) accept(visitor StatementVisitor)  error {
	return visitor.VisitBlockStatement(b)
}

func CreateBlock(statements []Statement) *BlockStatement {
	return &BlockStatement{
        Statements: statements,
	}
}
