/*  Copyright (c) 2013, Brian Hummer (brian@boggo.net)
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:
    * Redistributions of source code must retain the above copyright
      notice, this list of conditions and the following disclaimer.
    * Redistributions in binary form must reproduce the above copyright
      notice, this list of conditions and the following disclaimer in the
      documentation and/or other materials provided with the distribution.
    * Neither the name of the boggo.net nor the
      names of its contributors may be used to endorse or promote products
      derived from this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL BRIAN HUMMER BE LIABLE FOR ANY
DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
(INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

package main

import (
	"errors"
	//"fmt"
	"github.com/boggo/neat"
	"github.com/boggo/neat/decoder"
	"github.com/boggo/neat/popeval"
	"math"
)

type xorEval struct{}

func (eval xorEval) Evaluate(org *neat.Organism) (err error) {

	if org.Phenome == nil {
		err = errors.New("Cannot evaluate an org without a Phenome")
		org.Fitness = 0 // Minimal fitness
		return
	}

	in := [][]float64{{0, 0}, {0, 1}, {1, 0}, {1, 1}}
	out := []float64{0, 1, 1, 0}

	sum := float64(0)
	for i := 0; i < len(in); i++ {
		o, e2 := org.Analyze(in[i])
		if e2 != nil {
			err = e2
			org.Fitness = 0
			return
		}
		//fmt.Printf("Genome [%4d]: in %1.0f %1.0f out %1.0f guess was %7.6f err was %7.6f\n", org.ID, in[i][0], in[i][1], out[i], o[0], math.Pow(out[i]-o[0], 2))
		sum += math.Pow(out[i]-o[0], 2)
	}
	org.Fitness = float64(1) - math.Sqrt(sum/float64(4))

	return
}

func main() {

	neat.SetDecoder(decoder.NewNEAT())
	neat.SetPopEval(popeval.NewSerial())
	neat.SetOrgEval(xorEval{})
	neat.Iterate(100)

}
