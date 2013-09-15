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
)

type Species struct {
	ID          int           // Identifier for this species
	Orgs        OrganismSlice // Portion of the population belonging to this species
	Age         int           // Age of species
	BestFitness float64       // Best fitness this species has acheived
	BestFitAge  int           // Age when species achieved best fitness
	Example     *Organism     // Example organism for determining future members of this species
	currFitness float64       // The current generation's fitness
}

func (s Species) String() string {
	return fmt.Sprintf("Species [%5d] is %d old and has %d Organisms. Best fitness %6.4f was at %d",
		s.ID, s.Age, len(s.Orgs), s.BestFitness, s.BestFitAge)
}

func (s *Species) calcFitness() {
	sum := float64(0)
	for _, o := range s.Orgs {
		sum += o.Fitness
	}
	sum /= float64(len(s.Orgs))
	s.currFitness = sum

	if sum > s.BestFitness {
		s.BestFitness = sum
		s.BestFitAge = s.Age
	}
}

type SpeciesSlice []*Species

func (ss SpeciesSlice) Len() int           { return len(ss) }
func (ss SpeciesSlice) Swap(i, j int)      { ss[i], ss[j] = ss[j], ss[i] }
func (ss SpeciesSlice) Less(i, j int) bool { return ss[i].BestFitness < ss[j].BestFitness }

func (ss SpeciesSlice) Organisms() (orgs OrganismSlice) {
	orgs = make([]*Organism, 0, settings.PopulationSize)
	for _, s := range ss {
		for _, o := range s.Orgs {
			orgs = append(orgs, o)
		}
	}
	return
}
