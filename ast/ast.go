package ast

import (
	"bytes"
	"strings"

	"github.com/Savvelius/go-interp/token"
)

// interfaces ------------------

type Node interface {
	TokenLiteral() string // literal value of the token that node is associated with
	String() string       // source code representation of all data
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

// implementations -----------------

type Program struct {
	Statements []Statement
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

// structure: let + ident + assign + expression;
type LetStatement struct {
	Token token.Token // token.LET
	Name  *Identifier // name of variable this is binded to
	Value Expression  // value to be binded
}

func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}
	out.WriteByte(';')

	return out.String()
}
func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

type ReturnStatement struct {
	Token       token.Token // token.RETURN
	ReturnValue Expression  // expr to be returned
}

func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString("return ")
	// FIXME: delete branch after parsing exprs
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	out.WriteByte(';')

	return out.String()
}
func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }

// -x and !x
type PrefixExpression struct {
	Token    token.Token
	Operator string     // - or !
	Right    Expression // argument
}

func (pe *PrefixExpression) String() string {
	return "(" + pe.Operator + pe.Right.String() + ")"
}
func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }

type InfixExpression struct {
	Token    token.Token // operator token: +, -, *
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) String() string {
	return "(" + ie.Left.String() + " " + ie.Operator + " " + ie.Right.String() + ")"
}
func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }

type ExpressionStatement struct {
	Token      token.Token // first token of expression
	Expression Expression
}

func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}
func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }

type Identifier struct {
	Token token.Token // token.IDENT
	Value string      // variable name
}

func (i *Identifier) String() string       { return i.Value }
func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

type FunctionLiteral struct {
	Token      token.Token // token.FUNCTION
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, param := range fl.Parameters {
		params = append(params, param.String())
	}

	out.WriteString(fl.TokenLiteral())
	out.WriteByte('(')
	out.WriteString(strings.Join(params, ", "))
	out.WriteByte(')')
	out.WriteString(fl.Body.String())

	return out.String()
}

type IntegerLiteral struct {
	Token token.Token // token.INT
	Value int64
}

func (il *IntegerLiteral) String() string       { return il.Token.Literal }
func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string       { return b.Token.Literal }

type StringLiteral struct {
	Token token.Token // token.STRING
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return sl.Value }

type ArrayLiteral struct {
	Token    token.Token // [ token
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode()      {}
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }
func (al *ArrayLiteral) String() string {
	var out bytes.Buffer

	exprs := []string{}
	for _, expr := range al.Elements {
		exprs = append(exprs, expr.String())
	}
	out.WriteByte('[')
	out.WriteString(strings.Join(exprs, ", "))
	out.WriteByte(']')

	return out.String()
}

type HashLiteral struct {
	Token token.Token
	Pairs map[Expression]Expression
}

func (hl *HashLiteral) expressionNode()      {}
func (hl *HashLiteral) TokenLiteral() string { return hl.Token.Literal }
func (hl *HashLiteral) String() string {
	var out bytes.Buffer

	out.WriteByte('{')
	elems := []string{}
	for k, v := range hl.Pairs {
		elems = append(elems, k.String()+": "+v.String())
	}
	out.WriteString(strings.Join(elems, ", "))
	out.WriteByte('}')

	return out.String()
}

type IndexExpression struct {
	Token token.Token
	Index Expression
	Left  Expression
}

func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpression) String() string {
	var out bytes.Buffer

	out.WriteByte('(')
	out.WriteString(ie.Left.String())
	out.WriteByte('[')
	out.WriteString(ie.Index.String())
	out.WriteByte(']')
	out.WriteByte(')')

	return out.String()
}

type IfExpression struct {
	Token       token.Token // token.IF
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}

	return out.String()
}

type CallExpression struct {
	Token     token.Token // (
	Function  Expression  // Identifier or FunctionLiteral
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, arg := range ce.Arguments {
		args = append(args, arg.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteByte('(')
	out.WriteString(strings.Join(args, ", "))
	out.WriteByte(')')

	return out.String()
}

type BlockStatement struct {
	Token      token.Token // the { token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, stmt := range bs.Statements {
		out.WriteString(stmt.String())
	}

	return out.String()
}
