package main

import (
	"fmt"
)


type PrintVisitor struct{}

func (p *PrintVisitor) VisitLiteral(l *Literal) any {
	if l.Value == nil {
		return ""
	}
	return p.parenthesize(l.Value)
}

func (p *PrintVisitor) VisitUnary(u *Unary) any {
	return p.parenthesize(u.Operand.Lexeme, u.Right)
}

func (p *PrintVisitor) VisitBinary(b *Binary) any {
	return p.parenthesize(b.Operator.Lexeme, b.Left, b.Right)
}

func (p *PrintVisitor) VisitTernary(t *Ternary) any {
    return p.parenthesize("?:", t.Left, t.Center, t.Right)
}

func (p *PrintVisitor) VisitGrouping(g *Grouping) any {
	return p.parenthesize("group", g.Expression)
}

func (p *PrintVisitor) parenthesize(treeType interface{}, expressions ...Expression) string {
	str := ""
	str += "("
	str += fmt.Sprint(treeType)
	// for _, _ := range expressions {
	// 	str += " "
	// 	// str += expr.accept(p).(string)
	// }
	str += ")"
	return str
}
