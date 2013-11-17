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

type Loader interface {
	Load() (*Settings, error)
}

type Settings struct {

	// The number of invididuals in the population
	PopulationSize int

	// Size of the initial genome
	BiasCount   int
	InputCount  int
	OutputCount int

	// Coefficients for calculating distance between genomes
	ExcessCoefficient   float64
	DisjointCoefficient float64
	WeightCoefficient   float64

	// Probabilities for mutation
	MutateWeight        float64
	MutateWeightNew     float64
	MutateEnabled       float64
	MutateAddConnection float64
	MutateAddNode       float64
	MutateFuncType      float64
	MutateDelNode       float64 // Pruning phase
	MutateDelConnection float64 // Pruning phase
	PruneThreshold      float64 // Pruning phase threshold

	// Crossover and breeding probabilities
	Crossover          float64
	InterspeciesMating float64
	AgeToStagnation    int
	SurvivalPercent    float64 // Percent of a species to survive for mating
	EliteCount         int     // Number within a species to survive into the next generation
	CompatThreshold    float64 // Compatiblity threshold for adding a genome to a species

	// Runtime settings
	ArchiveFrequency int // Frequency to archive the population. 0 = archive every iteration
	ReportFrequency  int // Frequency to report on the population. 0 = report every iteration
}
