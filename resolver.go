package main

type Resolver struct {
	*Interpreter
	*Lox
	Scopes []map[string]bool
}

func CreateResolver(interpreter *Interpreter, lox *Lox) *Resolver {
	return &Resolver{
		Interpreter: interpreter,
		Scopes:      make([]map[string]bool, 0),
		Lox:         lox,
	}
}

func (r *Resolver) resolve(statements []Statement) {
	for _, stmt := range statements {
		// Check this :
		stmt.accept(r)
	}
}

func (r *Resolver) VisitVarDeclaration(v *VarDeclaration) (any, error) {
	r.declare(v.Identifier)
	r.resolveExpr(v.Expr)
	r.define(v.Identifier)

	return nil, nil
}

func (r *Resolver) VisitFunctionDeclaration(f *FunctionDeclaration) (any, error) {
	r.declare(f.Identifier)
	r.define(f.Identifier)

	r.beginScope()
	defer r.endScope()

	for _, param := range f.Params {
		r.declare(param)
		r.define(param)
	}

	for _, stmt := range f.Stmts {
		r.resolveStmt(stmt)
	}

	return nil, nil
}

func (r *Resolver) VisitReturnStatement(ret *ReturnStatement) (any, error) {
	if ret.Expr != nil {
		r.resolveExpr(ret.Expr)
	}
	return nil, nil
}

func (r *Resolver) VisitIfStatement(i *IfStatement) (any, error) {
	r.resolveExpr(i.Expr)

	r.resolveStmt(i.IfStmt)
	if i.ElseStmt != nil {
		r.resolveStmt(i.ElseStmt)
	}

	return nil, nil
}

func (r *Resolver) VisitWhileStatement(w *WhileStatement) (any, error) {
	r.resolveExpr(w.Expr)
	r.resolveStmt(w.Stmt)

	return nil, nil
}

func (r *Resolver) VisitBlockStatement(b *BlockStatement) (any, error) {
	r.beginScope()
	defer r.endScope()

	for _, stmt := range b.Statements {
		r.resolveStmt(stmt)
	}

	return nil, nil

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
	r.resolveExpr(f.Identifier)

	for _, arg := range *f.Args {
		r.resolveExpr(arg)
	}
	return nil, nil
}

func (r *Resolver) VisitVarAssignment(v *VarAssignment) (any, error) {
	r.resolveFinal(v.Token, v)
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
	r.resolveFinal(i.name, i)
	return nil, nil
}

func (r *Resolver) VisitLiteral(l *Literal) any {
	return nil
}

func (r *Resolver) resolveFinal(token *Token, expr Expression) {
	for idx := len(r.Scopes) - 1; idx >= 0; idx -= 1 {
		curr := r.Scopes[idx]
		val, ok := curr[token.Lexeme]
		if ok {
			if !val {
				r.HadError = true
				return
			}
			dist := len(r.Scopes) - 1 - idx
			r.Interpreter.Locals[expr] = dist

			break
			// panic(fmt.Sprintf("Get the distance : %d\n", dist))
		}
	}
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
