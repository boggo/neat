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
	"sync"
)

type Decoder interface {
	Decode(genome *Genome) (phenome Phenome, err error)
}

type OrgEval interface {
	Evaluate(org *Organism) (err error)
}

type PopEval interface {
	Evaluate(pop *Population, orgEval OrgEval) (err error)
}

func Iterate(settings *Settings, n int, dcode Decoder, popEval PopEval, orgEval OrgEval, arch Archiver, rep Reporter) {

	var err error

	// Phase search parameters
	var pth float64 // Pruning threshold
	var cmplx bool  // Switch between complexifying (true) and simplifying (false)
	var addNode, delNode, addConn, delConn, cross float64 // original values
	addNode = settings.MutateAddNode
	delNode = settings.MutateDelNode
	addConn = settings.MutateAddConnection
	delConn = settings.MutateDelConnection
	cross = settings.Crossover
	cmplx = true	// Start with complexifying

	// Restore the population
	var population *Population
	if arch != nil {
		population, err = arch.Restore()
		if err != nil {
			fmt.Println("Restore failed:", err) // Will begin a new population
		} else {
			pth = population.MPC() + settings.PruneThreshold
		}
	}

	// Create the innovation tracker
	inno := newInnovation(population)
	defer inno.close()

	//Iterate
	for i := 0; i < n; i++ {

		// Ensure the current population
		if population == nil {
			population, err = initialPopulation(settings, inno)
			pth = population.MPC() + settings.PruneThreshold
		} else {

			// Determine if the search should switch to decomplexifying
			if cmplx {}
			mpc := population.MPC()
			cmplx = (settings.PruneThreshold > 0) && (mpc < pth)
			} else {
				m = population.MPC()
				if nochg > settings.PruneFloor {
					cmplx = true
				}
				
			}
			if cmplx {
				settings.MutateAddNode = addNode
				settings.MutateAddConnection = addConn
				settings.MutateDelNode = 0
				settings.MutateDelConnection = 0
				settings.Crossover = cross
			} else {
				settings.MutateAddNode = 0
				settings.MutateAddConnection = 0
				settings.MutateDelNode = delNode
				settings.MutateDelConnection = delConn
				settings.Crossover = 0
			}
			// Roll to the next generation
			population, err = rollPop(settings, inno, population)
		}
		if err != nil {
			panic(err)
		}

		// Ensure every organism is decoded
		var w sync.WaitGroup
		orgs := population.Species.Organisms(settings)
		for _, o := range orgs {
			if o.Phenome == nil {
				w.Add(1)
				go func(o *Organism) {
					o.Phenome, err = dcode.Decode(o.Genome)
					if err != nil {
						// Do what exactly?
					}
					w.Done()
				}(o)
			}
		}
		w.Wait()

		// Evaluate each organism
		err = popEval.Evaluate(population, orgEval)
		if err != nil {
			panic(err)
		}

		// Archive the population
		if arch != nil && (i == n-1 ||
			(settings.ArchiveFrequency == 0 || i%settings.ArchiveFrequency == 0)) {
			err = arch.Archive(population)
			if err != nil {
				panic(err)
			}
		}

		// Report the population
		if rep != nil && (i == n-1 || (settings.ReportFrequency == 0 || i%settings.ReportFrequency == 0)) {
			err = rep.Report(population)
			if err != nil {
				panic(err)
			}
		}

	}

}
