/*
CREDIT WHERE CREDIT IS DUE:
This experiment is based on the XOR experiement in neat-python
(https://code.google.com/p/neat-python/), specifically the Evaluate
function below is based on the eval_fitness function in xor2.py.
That code is GPL3 and copyright belongs to that software's authors.

The remaining code in this experiment was written by Brian Hummer (brian@boggo.net)
and is released under GPL3 because the example used was GPL3. The libraries used,
though, are 3-clause, "new" BSD licensed.

GPL3 license:

Copyright (C) 2013 Brian Hummer (brian@boggo.net)

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package main

import (
	"errors"
	"github.com/boggo/neat"
	"github.com/boggo/neat/archiver"
	"github.com/boggo/neat/decoder"
	"github.com/boggo/neat/popeval"
	"github.com/boggo/neat/reporter"
	"github.com/boggo/neat/settings"
	"math"
)

/* Original eval function from neat-python:

def eval_fitness(population):
    for chromo in population:
        net = nn.create_ffphenotype(chromo)

        error = 0.0
        #error_stanley = 0.0
        for i, inputs in enumerate(INPUTS):
            net.flush() # not strictly necessary in feedforward nets
            output = net.sactivate(inputs) # serial activation
            error += (output[0] - OUTPUTS[i])**2

            #error_stanley += math.fabs(output[0] - OUTPUTS[i])

        #chromo.fitness = (4.0 - error_stanley)**2 # (Stanley p. 43)
        chromo.fitness = 1 - math.sqrt(error/len(OUTPUTS))

*/

var (
	INPUTS  [][]float64
	OUTPUTS []float64
)

func init() {
	INPUTS = [][]float64{{0, 0}, {0, 1}, {1, 0}, {1, 1}}
	OUTPUTS = []float64{0, 1, 1, 0}
}

type xorEval struct{}

func (eval xorEval) Evaluate(org *neat.Organism) (err error) {

	if org.Phenome == nil {
		err = errors.New("Cannot evaluate an org without a Phenome")
		org.Fitness = []float64{0} // Minimal fitness
		return
	}

	e := float64(0)
	for i, inputs := range INPUTS {
		output, err2 := org.Analyze(inputs)
		if err2 != nil {
			err = err2
			org.Fitness = []float64{0}
			return
		}
		e += (output[0] - OUTPUTS[i]) * (output[0] - OUTPUTS[i])
	}
	org.Fitness = []float64{float64(1) - math.Sqrt(e/float64(len(OUTPUTS)))}

	return
}

func main() {

	// Load the settings
	ldr := settings.NewJSON("xor-settings.json")
	s, err := ldr.Load()
	if err != nil {
		panic(err)
	}

	// Create the archiver
	a := archiver.NewJSON("xor-pop.json")

	// Create the reporter
	r := reporter.NewConsole()

	// Create the evaluators
	o := &xorEval{}
	p := popeval.NewConcurrent()

	// Create the decoder
	d := decoder.NewNEAT()

	// Iterate the experiment
	neat.Iterate(s, 100, d, p, o, a, r)
}
