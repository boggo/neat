package neat

import (
"github.com/boggo/neural"
)
// A Phenome takes an slice of inputs and generates the corresponding
// outputs based on its internal "black box"
type Phenome interface {
    Analyze(inputs []float64) (outputs []float64, err error)
}

type networkPhenome struct {
	network *neural.Network
}

func (p *Phenome) Analyze(inputs []float64) (outputs []float64, err error) {
	outputs = p.Activate(inputs)
	return
}