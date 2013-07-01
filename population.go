package neat


type Population struct {

    // Generation number
    Generation uint32

    // Slice of species - this is what we track
    Species SpeciesSlice

    // Slice of organism - this is what we evaluate
    orgs OrganismSlice
}

type PopulationEvaluator interface {
    Evaluate(pop *Population, orgEval OrganismEvaluator) (err error)
}

func NewPopEval() PopulationEvaluator {
    return &serialPopEval{}
}

// Default Population Evaluators. Processes organisms one after the other.
// This should be replaced with a concurrent (or even more advanced) evaluator
// for production
type serialPopEval struct{}

// Evaluates the entire population
func (pe serialPopEval) Evaluate(pop *Population, orgEval OrganismEvaluator) (err error) {

    // Iterate the organisms in the population
    // TODO: make a Go routine out of this so it can execute concurrently
    for i, _ := range pop.orgs {
        pop.orgs[i].Fitness, err = orgEval.Evaluate(pop.orgs[i])
    }

    // Evaluation complete
    return
}
