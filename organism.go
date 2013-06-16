package neat

// An Organism encompasses both the encoding, Genome, and the expresssion
// of that encoding, Phenome.
type Organism struct {
    Phenome // Provides the functionality to evaluate inputs
    *Genome // Provides the ID and Fitness attributes
}

// Slice of Organisms
type OrganismSlice []*Organism

// Len is the number of elements in the collection.
func (orgs OrganismSlice) Len() int { return len(orgs) }

// Less returns whether the element with index i should sort
// before the element with index j.
func (orgs OrganismSlice) Less(i, j int) bool {
    // Compare the organisms' primary fitness values
    return orgs[i].Fitness[0] < orgs[j].Fitness[0]
}

// Swap swaps the elements with indexes i and j.
func (orgs OrganismSlice) Swap(i, j int) {
    orgs[i], orgs[j] = orgs[j], orgs[i]
}

// Evaluates an Organism, returning its fitness
type OrganismEvaluator interface {
	Evaluate(org *Organism) (fitness []float64, err error)
}