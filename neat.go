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

package neat

import (
	"fmt"
	"math"
	"sort"
)

var (
	population *Population // The current population
	decoder    Decoder     // Decodes a genome into a phenome
	orgEval    OrgEval     // Evaluates a single organism
	popEval    PopEval     // Evaluates the whole population
)

func SetDecoder(d Decoder) {
	decoder = d
}

func SetOrgEval(o OrgEval) {
	orgEval = o
}

func SetPopEval(p PopEval) {
	popEval = p
}

type Decoder interface {
	Decode(genome *Genome) (phenome Phenome, err error)
}

type OrgEval interface {
	Evaluate(org *Organism) (err error)
}

type PopEval interface {
	Evaluate(pop *Population, orgEval OrgEval) (err error)
}

func Iterate(iters int) {

	var err error

	for i := 0; i < iters; i++ {

		// Ensure the current population
		if population == nil {
			population, err = initialPopulation()
		} else {
			err = rollPop()
		}
		if err != nil {
			panic(err)
		}

		// Ensure every organism is decoded
		orgs := population.Species.Organisms()
		for _, o := range orgs {
			if o.Phenome == nil {
				o.Phenome, err = decoder.Decode(o.Genome)
				if err != nil {
					// Do what exactly?
				}
			}
		}

		// Evaluate each organism
		err = popEval.Evaluate(population, orgEval)
		if err != nil {
			panic(err)
		}

		//DumpPopulation()
		Best()
	}

	//DumpPopulation()
	Best()
}

// Debug
func DumpPopulation() {
	fmt.Println("Current Population:")
	fmt.Println(population)
	for _, s := range population.Species {
		fmt.Printf("\t%v\n", s)
		for _, o := range s.Orgs {
			DumpOrg(o)
		}
	}
}

func DumpOrg(o *Organism) {
	fmt.Printf("\t\t%v\n", o.Genome)
	for _, ng := range o.Nodes {
		fmt.Printf("\t\t\t%v\n", ng)
	}
	for _, cg := range o.Conns {
		fmt.Printf("\t\t\t%v\n", cg)
	}
}
func Best() {

	fmt.Println("Current Population:")
	fmt.Println(population)
	for _, s := range population.Species {
		fmt.Printf("\t%v\n", s)
	}
	orgs := population.Species.Organisms()
	sort.Sort(sort.Reverse(orgs))
	org := orgs[0]
	fmt.Println("Best org", org)
	for _, ng := range org.Nodes {
		fmt.Printf("\t\t\t%v\n", ng)
	}
	for _, cg := range org.Conns {
		fmt.Printf("\t\t\t%v\n", cg)
	}
	in := [][]float64{{0, 0}, {0, 1}, {1, 0}, {1, 1}}
	out := []float64{0, 1, 1, 0}

	sum := float64(0)

	for i := 0; i < len(in); i++ {
		o, _ := org.Analyze(in[i])
		e := math.Pow(out[i]-o[0], 2)
		sum += e
		fmt.Printf("%f XOR %f = %f, output was %8.6f error was %8.6f\n", in[i][0], in[i][1], out[i], o[0], e)
	}
	fmt.Println("Org fitness is therefore ", float64(1)-math.Sqrt(sum/float64(4)))
	if org.Fitness > 0.9 {
		panic("No error just done!")
	}
}
