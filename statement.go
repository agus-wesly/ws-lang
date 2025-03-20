package main

import (
	"errors"
	"fmt"
)

type Statement interface {
	accept(StatementVisitor) (any, error)
}

type StatementVisitor interface {
	VisitExpressionStatement(e *ExpressionStatement) (any, error)
	VisitPrintStatement(p *PrintStatement) error
	VisitVarDeclaration(v *VarDeclaration) (any, error)
	VisitBlockStatement(v *BlockStatement) (any, error)
	VisitIfStatement(i *IfStatement) (any, error)
	VisitWhileStatement(w *WhileStatement) (any, error)
	VisitFunctionDeclaration(f *FunctionDeclaration) (any, error)
	VisitReturnStatement(r *ReturnStatement) (any, error)
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
	Identifier *Token
	Expr       Expression
}

func (v *VarDeclaration) accept(visitor StatementVisitor) (any, error) {
	return visitor.VisitVarDeclaration(v)
}

func CreateVarDeclaration(expr Expression, identifier *Token) *VarDeclaration {
	return &VarDeclaration{
		Expr:       expr,
		Identifier: identifier,
	}
}

type FunctionDeclaration struct {
	Identifier *Token
	Params     []*Token
	Stmts      []Statement
}

func (f *FunctionDeclaration) accept(visitor StatementVisitor) (any, error) {
	return visitor.VisitFunctionDeclaration(f)
}

func CreateFunctionDeclaration(identifier *Token, params []*Token, stmts []Statement) *FunctionDeclaration {
	return &FunctionDeclaration{
		Identifier: identifier,
		Params:     params,
		Stmts:      stmts,
	}
}

type VarAssignment struct {
	Token *Token
	Expr  Expression
}

func (v *VarAssignment) accept(visitor ExpressionVisitor) (any, error) {
	return visitor.VisitVarAssignment(v)
}

func CreateVarAssignment(token *Token, expr Expression) *VarAssignment {
	return &VarAssignment{
		Token: token,
		Expr:  expr,
	}
}

type BlockStatement struct {
	Statements []Statement
}

func (b *BlockStatement) accept(visitor StatementVisitor) (any, error) {
	return visitor.VisitBlockStatement(b)
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
	return visitor.VisitIfStatement(i)
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
	Stmt []Statement
}

func (w *WhileStatement) accept(visitor StatementVisitor) (any, error) {
	return visitor.VisitWhileStatement(w)
}

func CreateWhileStatement(expr Expression, stmt []Statement) *WhileStatement {
	return &WhileStatement{
		Expr: expr,
		Stmt: stmt,
	}
}

type BreakStatement struct{}

var BreakStmtErr = errors.New("BreakStatement")

func (b *BreakStatement) accept(visitor StatementVisitor) (any, error) {
	return nil, BreakStmtErr
}

func CreateBreakStatement() *BreakStatement {
	return &BreakStatement{}
}

type ReturnStatement struct {
	Expr Expression
	*Token
}

func (r *ReturnStatement) Error() string {
	return fmt.Sprintf("[line %d] Runtime Error : Illegal return statement\n", r.Token.Line)
}

func (r *ReturnStatement) accept(visitor StatementVisitor) (any, error) {
	return visitor.VisitReturnStatement(r)
}

func CreateReturnStatement(token *Token, expr Expression) *ReturnStatement {
	return &ReturnStatement{
		Expr:  expr,
		Token: token,
	}
}
