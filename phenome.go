package neat

import (
    "github.com/boggo/neural"
)

// A Phenome takes an slice of inputs and generates the corresponding
// outputs based on its internal "black box"
type Phenome interface {
    Analyze(inputs []float64) (outputs []float64, err error)
}

// Default implementation of Phenome uses a neural.Network to analyze the
// inputs
type networkPhenome struct {
    network *neural.Network
}

// Analyzes the inputs and returns the results as a slice of float64. This
// implementation is for the networkPhenome. Network does not return an error
// in Activate() so the error value is nil upon return.
func (p *networkPhenome) Analyze(inputs []float64) (outputs []float64, err error) {
    outputs = p.network.Activate(inputs)
    return
}
