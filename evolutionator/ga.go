package evolutionator

import (
	"fmt"
	"math/rand"

	"github.com/MaxHalford/eaopt"
	"github.com/robertkrimen/otto/ast"
)

// MakeJSGenome ignores the random initialization, all individuals are
// clones of the first parsed AST
func MakeJSGenome(ast *ast.Program) func(rng *rand.Rand) eaopt.Genome {
	return func(rng *rand.Rand) eaopt.Genome {
		g := jsGenome{ast}
		return g.Clone()
	}
}

func Optimize(ast *ast.Program) (*ast.Program, error) {
	var conf = eaopt.NewDefaultGAConfig()
	conf.NPops = 1
	var ga, err = conf.NewGA()
	if err != nil {
		return nil, err
	}

	// Add a custom print function to track progress
	ga.Callback = func(ga *eaopt.GA) {
		fmt.Printf("Best fitness at generation %d: %f\n", ga.Generations, ga.HallOfFame[0].Fitness)
	}

	// Run the GA
	err = ga.Minimize(MakeJSGenome(ast))
	if err != nil {
		return nil, err
	}

	// Best genome
	g := ga.HallOfFame[0].Genome

	jsGenome := g.(jsGenome)
	return jsGenome.ast, nil
}
