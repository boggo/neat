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

package reporter

import (
	"github.com/boggo/neat"
	"math"
)

func max(x []float64) float64 {
	m := float64(0)
	for _, a := range x {
		if a > m {
			m = a
		}
	}
	return m
}

func min(x []float64) float64 {
	m := math.MaxFloat64
	for _, a := range x {
		if a < m {
			m = a
		}
	}
	return m
}

func avg(x []float64) float64 {
	n := float64(len(x))
	if n == 0 {
		return 0
	}

	s := float64(0)
	for _, a := range x {
		s += a
	}

	return s / n
}

func avgFit(s *neat.Species) float64 {
	x := make([]float64, len(s.Orgs))
	for i, o := range s.Orgs {
		x[i] = o.Fitness[0]
	}
	return avg(x)
}

func avgNod(s *neat.Species) float64 {
	x := make([]float64, len(s.Orgs))
	for i, o := range s.Orgs {
		x[i] = float64(len(o.Nodes))
	}
	return avg(x)
}

func avgCon(s *neat.Species) float64 {
	x := make([]float64, len(s.Orgs))
	for i, o := range s.Orgs {
		x[i] = float64(len(o.Conns))
	}
	return avg(x)
}

type speciesSort struct {
	species []*neat.Species
}

func (ss *speciesSort) Len() int           { return len(ss.species) }
func (ss *speciesSort) Swap(i, j int)      { ss.species[i], ss.species[j] = ss.species[j], ss.species[i] }
func (ss *speciesSort) Less(i, j int) bool { return ss.species[i].ID < ss.species[j].ID }
