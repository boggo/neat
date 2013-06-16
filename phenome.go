package neat

// A Phenome takes an slice of inputs and generates the corresponding
// outputs based on its internal "black box"
type Phenome interface {
    Activate(inputs []float64) (outputs []float64)
}

