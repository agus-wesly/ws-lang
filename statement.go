package main

type Statement interface {
	accept(StatementVisitor) (any, error)
}

type StatementVisitor interface {
	VisitExpressionStatement(e *ExpressionStatement) (any, error)
	VisitPrintStatement(p *PrintStatement) error
	VisitVarDeclaration(v *VarDeclaration) (any, error)
	VisitBlockStatement(v *BlockStatement) error
	VisitIfStatement(i *IfStatement) error
	VisitWhileStatement(w *WhileStatement) error
}

type ExpressionStatement struct {
	Expr Expression
}

func (e *ExpressionStatement) accept(visitor StatementVisitor) (any, error) {
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

func (p *PrintStatement) accept(visitor StatementVisitor) (any, error) {
	return nil, visitor.VisitPrintStatement(p)
}

func CreatePrintStatement(expr Expression) *PrintStatement {
	return &PrintStatement{
		Expr: expr,
	}
}

type VarDeclaration struct {
	Identifier Token
	Expr       Expression
}

func (v *VarDeclaration) accept(visitor StatementVisitor) (any, error) {
	return visitor.VisitVarDeclaration(v)
}

func CreateVarDeclaration(expr Expression, identifier Token) *VarDeclaration {
	return &VarDeclaration{
		Expr:       expr,
		Identifier: identifier,
	}
}

type VarAssignment struct {
	Token Token
	Expr  Expression
}

func (v *VarAssignment) accept(visitor ExpressionVisitor) (any, error) {
	return visitor.VisitVarAssignment(v)
}

func CreateVarAssignment(token Token, expr Expression) *VarAssignment {
	return &VarAssignment{
		Token: token,
		Expr:  expr,
	}
}

type BlockStatement struct {
	Statements []Statement
}

func (b *BlockStatement) accept(visitor StatementVisitor) (any, error) {
	return nil, visitor.VisitBlockStatement(b)
}

func CreateBlock(statements []Statement) *BlockStatement {
	return &BlockStatement{
		Statements: statements,
	}
}

type IfStatement struct {
	Expr     Expression
	IfStmt   Statement
	ElseStmt Statement
}

func (i *IfStatement) accept(visitor StatementVisitor) (any, error) {
	return nil, visitor.VisitIfStatement(i)
}

func CreateIfStatement(expr Expression, ifStmt Statement, elseStmt Statement) *IfStatement {
	return &IfStatement{
		Expr:     expr,
		IfStmt:   ifStmt,
		ElseStmt: elseStmt,
	}
}

type WhileStatement struct {
	Expr Expression
	Stmt Statement
}

func (w *WhileStatement) accept(visitor StatementVisitor) (any, error) {
	return nil, visitor.VisitWhileStatement(w)
}

func CreateWhileStatement(expr Expression, Stmt Statement) *WhileStatement {
	return &WhileStatement{
		Expr: expr,
		Stmt: Stmt,
	}
}
