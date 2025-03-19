package main

import "fmt"

type Resolver struct {
	*Interpreter
	Scopes []map[string]bool
}

func (r *Resolver) VisitVarDeclaration(v *VarDeclaration) (any, error) {
	// Create and push to scope
	r.beginScope()
	defer r.endScope()

	r.declare(v.Identifier)
	r.resolveExpr(v.Expr)
	r.define(v.Identifier)

	return nil, nil
}

func (r *Resolver) VisitFunctionDeclaration(f *FunctionDeclaration) (any, error) {
	r.define(f.Identifier)
	r.declare(f.Identifier)

	r.beginScope()
	defer r.endScope()

	for _, param := range f.Params {
		r.define(param)
		r.declare(param)
	}

	for _, stmt := range f.Stmts {
		r.resolveStmt(stmt)
	}

	return nil, nil
}

func (r *Resolver) VisitIfStatement(i *IfStatement) error {
	r.resolveExpr(i.Expr)

	r.beginScope()
	defer r.endScope()

	r.resolveStmt(i.IfStmt)
	if i.ElseStmt != nil {
		r.resolveStmt(i.ElseStmt)
	}

	return nil
}

func (r *Resolver) VisitWhileStatement(w *WhileStatement) error {
	r.resolveExpr(w.Expr)

	r.beginScope()
	defer r.endScope()

	r.resolveStmt(w.Stmt)

	return nil
}

func (r *Resolver) VisitBlockStatement(b *BlockStatement) error {
	r.beginScope()
	defer r.endScope()

	for _, stmt := range b.Statements {
		r.resolveStmt(stmt)
	}

	return nil

}

func (r *Resolver) VisitPrintStatement(p *PrintStatement) error {
	r.resolveExpr(p.Expr)
	return nil
}

func (r *Resolver) VisitExpressionStatement(e *ExpressionStatement) (any, error) {
	r.resolveExpr(e.Expr)
	return nil, nil
}

func (r *Resolver) VisitFunction(f *Function) (any, error) {
	// TODO : change in Function struct to use Identifier not expr
	r.resolveFinal(f.Identifier)
	return nil, nil
}

func (r *Resolver) VisitVarAssignment(v *VarAssignment) (any, error) {
	r.resolveFinal(v.Token)
	r.resolveExpr(v.Expr)

	return nil, nil
}

func (r *Resolver) VisitGrouping(g *Grouping) (any, error) {
	r.resolveExpr(g.Expression)

	return nil, nil
}

func (r *Resolver) VisitLogicalOperator(l *LogicalOperator) (any, error) {
	r.resolveExpr(l.Left)
	r.resolveExpr(l.Right)

	return nil, nil
}

func (r *Resolver) VisitTernary(t *Ternary) (any, error) {
	r.resolveExpr(t.Left)
	r.resolveExpr(t.Center)
	r.resolveExpr(t.Right)

	return nil, nil
}

func (r *Resolver) VisitBinary(b *Binary) (any, error) {
	r.resolveExpr(b.Left)
	r.resolveExpr(b.Right)

	return nil, nil
}

func (r *Resolver) VisitUnary(u *Unary) (any, error) {
	r.resolveExpr(u.Right)

	return nil, nil
}

func (r *Resolver) VisitIdentifier(i *Identifier) (any, error) {
	// Begin search
	r.resolveFinal(i.name)
	// Attach to the node
	return nil, nil
}

func (r *Resolver) VisitLiteral(l *Literal) any {
	return nil
}

func (r *Resolver) resolveFinal(token *Token) {
	for idx := len(r.Scopes) - 1; idx >= 0; idx -= 1 {
		curr := r.Scopes[idx]
		val, ok := curr[token.Lexeme]
		if ok {
			if !val {
				panic("Found but not yet defined")
			}
			panic(fmt.Sprintf("Get the distance : %d\n", len(r.Scopes)-1-idx))
		}
	}

	panic("TODO : search in global env")
}

func (r *Resolver) declare(token *Token) {
	if r.isEmpty() {
		return
	}

	cur := r.Scopes[len(r.Scopes)-1]
	cur[token.Lexeme] = false
}

func (r *Resolver) define(token *Token) {
	if r.isEmpty() {
		return
	}

	cur := r.Scopes[len(r.Scopes)-1]
	_, ok := cur[token.Lexeme]
	if ok {
		cur[token.Lexeme] = true
	}
}

func (r *Resolver) beginScope() {
	newScope := map[string]bool{}
	r.Scopes = append(r.Scopes, newScope)
}

func (r *Resolver) endScope() {
	r.Scopes = (r.Scopes[:len(r.Scopes)-1])
}

func (r *Resolver) resolveExpr(expr Expression) {
	expr.accept(r)
}

func (r *Resolver) resolveStmt(stmt Statement) {
	stmt.accept(r)
}

func (r *Resolver) isEmpty() bool {
	return len(r.Scopes) == 0
}
