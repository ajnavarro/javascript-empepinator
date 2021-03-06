package evolutionator

import (
	"math/rand"

	"github.com/robertkrimen/otto/ast"
)

type Mutator struct {
	rand *rand.Rand
	ids  []*ast.Identifier
}

func NewMutator(rand *rand.Rand) *Mutator {
	return &Mutator{
		rand: rand,
	}
}

func (e *Mutator) Mutate(a *ast.Program) *ast.Program {
	if a == nil {
		return nil
	}

	for i, b := range a.Body {
		a.Body[i] = e.evalStatement(b)
	}

	return a
}

func (e *Mutator) evalStatement(stat ast.Statement) ast.Statement {
	switch t := stat.(type) {
	case *ast.FunctionStatement:
		if e.hit(3) {
			t.Function.Body = e.evalStatement(t.Function.Body)
		}
	case *ast.BlockStatement:
		var outList []ast.Statement
		for _, s := range t.List {
			if !e.hit(5) {
				outList = append(outList, e.evalStatement(s))
			}
		}
		t.List = outList
	case *ast.ReturnStatement:
		if e.hit(3) {
			t.Argument = e.evalExpression(t.Argument)
		}
	case *ast.VariableStatement:
		var outList []ast.Expression
		for _, ex := range t.List {
			if !e.hit(5) {
				outList = append(outList, e.evalExpression(ex))
			}
		}

		t.List = outList
	default:
		println(t)
		panic("STATEMENT NOT FOUND")
	}

	return stat
}

func (e *Mutator) evalExpression(expr ast.Expression) ast.Expression {
	switch t := expr.(type) {
	case *ast.BinaryExpression:
		if e.hit(3) {
			t.Left = e.evalExpression(t.Left)
		}

		if e.hit(3) {
			t.Right = e.evalExpression(t.Right)
		}
	case *ast.UnaryExpression:
		if e.hit(3) {
			t.Operand = e.evalExpression(t.Operand)
		}
	case *ast.VariableExpression:
		t.Initializer = e.evalExpression(t.Initializer)
	case *ast.Identifier:
		if e.hit(3) {
			switch len(e.ids) {
			case 1:
				expr = e.ids[0]
			case 0:
				break
			default:
				expr = e.ids[e.rand.Intn(len(e.ids)-1)]
			}
		}

		e.ids = append(e.ids, t)
	default:
		panic(t)
	}

	return expr
}

func (e *Mutator) hit(max int32) bool {
	return e.rand.Int31n(max) == 1
}
