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
	"encoding/json"
	"fmt"
	"github.com/boggo/neural"
	"strconv"
)

var (
	marker int
	id     int
)

// Encoding for the node in the neural network
type NodeGene struct {
	Marker int             // Innovation marker for this gene
	Type   neural.NodeType // Network node type
	X, Y   float64         // 2-D Position of this node within the network
}

func (ng NodeGene) String() string {
	var t string
	switch ng.Type {
	case neural.BIAS:
		t = "BIAS"
	case neural.INPUT:
		t = "INPUT"
	case neural.OUTPUT:
		t = "OUTPUT"
	case neural.HIDDEN:
		t = "HIDDEN"
	default:
		t = "UNKNOWN"
	}
	return fmt.Sprintf("NodeGene [%4d] %7v at %3.2f, %3.2f", ng.Marker, t, ng.X, ng.Y)
}

type NodeGeneMap map[int]*NodeGene

func (im NodeGeneMap) MarshalJSON() (bytes []byte, err error) {

	// Convert to a string-keyed map
	sk := make(map[string]*NodeGene)

	// Copy the values
	for k, v := range im {
		ks := strconv.Itoa(k)
		sk[ks] = v
	}

	// Marshal the string-keyed map
	bytes, err = json.Marshal(sk)
	return
}

func (im *NodeGeneMap) UnmarshalJSON(bytes []byte) (err error) {

	// Unmarshal the string-keyed map
	sk := make(map[string]*NodeGene)
	err = json.Unmarshal(bytes, &sk)
	if err != nil {
		return
	}

	// Copy the values
	*im = make(map[int]*NodeGene)
	var ki int
	for k, v := range sk {
		ki, err = strconv.Atoi(k)
		if err != nil {
			return
		}
		(*im)[ki] = v
	}

	return

}

func cloneNode(source *NodeGene) (clone *NodeGene) {
	clone = &NodeGene{Marker: source.Marker, Type: source.Type, X: source.X, Y: source.Y}
	return
}

// Encoding for the connection in the neural network
type ConnGene struct {
	Marker         int     // Innovation marker for this gene
	Source, Target int     // Innovation markers for the source and target node genes
	Weight         float64 // Weight applied during activation
	Enabled        bool    // Is this connection gene enabled?
}

type ConnGeneMap map[int]*ConnGene

func (im ConnGeneMap) MarshalJSON() (bytes []byte, err error) {

	// Convert to a string-keyed map
	sk := make(map[string]*ConnGene)

	// Copy the values
	for k, v := range im {
		ks := strconv.Itoa(k)
		sk[ks] = v
	}

	// Marshal the string-keyed map
	bytes, err = json.Marshal(sk)
	return
}

func (im *ConnGeneMap) UnmarshalJSON(bytes []byte) (err error) {

	// Unmarshal the string-keyed map
	sk := make(map[string]*ConnGene)
	err = json.Unmarshal(bytes, &sk)
	if err != nil {
		return
	}

	// Copy the values
	*im = make(map[int]*ConnGene)
	var ki int
	for k, v := range sk {
		ki, err = strconv.Atoi(k)
		if err != nil {
			return
		}
		(*im)[ki] = v
	}

	return

}

func (cg ConnGene) String() string {
	var e string
	if cg.Enabled {
		e = "ENABLED"
	} else {
		e = "DISABLED"
	}
	return fmt.Sprintf("ConnGene [%4d] Src: %4d Tgt: %4d Wgt: %+8.6f %v", cg.Marker, cg.Source,
		cg.Target, cg.Weight, e)
}

func cloneConn(source *ConnGene) (clone *ConnGene) {
	clone = &ConnGene{Marker: source.Marker, Source: source.Source, Target: source.Target,
		Weight: source.Weight, Enabled: source.Enabled}
	return
}

// Encoded organism
type Genome struct {
	ID      int         // Public identifier for this genome
	Nodes   NodeGeneMap // Collection of node genes identified by their markers
	Conns   ConnGeneMap // Collection of conn genes identified by their markers
	Fitness []float64   // Fitness of this Genome
}

// Describes the genome
func (g Genome) String() string {
	return fmt.Sprintf("Genome [%4d] has %3d Nodes and %3d Conns, Fitness: %8.4f", g.ID,
		len(g.Nodes), len(g.Conns), g.Fitness)
}

// Creates a deep copy of the genome
func cloneGenome(source *Genome, id int) (clone *Genome) {
	clone = &Genome{ID: id, Fitness: source.Fitness,
		Nodes: make(map[int]*NodeGene), Conns: make(map[int]*ConnGene)}
	for k, v := range source.Nodes {
		clone.Nodes[k] = cloneNode(v)
	}
	for k, v := range source.Conns {
		clone.Conns[k] = cloneConn(v)
	}
	return clone
}

// Creates the initial genome to seed the population
func initialGenome(settings *Settings, inno *innovation) (genome *Genome, err error) {

	// Shortcuts to settings values
	biasCount := settings.BiasCount
	inputCount := settings.InputCount
	outputCount := settings.OutputCount

	// Create the new genome with a negative (i.e., invalid) ID
	genome = &Genome{ID: -1, Nodes: make(map[int]*NodeGene), Conns: make(map[int]*ConnGene)}

	// Construct the nodes
	// Create the bias and input nodes
	var step float64
	var ng *NodeGene
	step = 0
	if biasCount+inputCount > 1 {
		step = 1.0 / float64(biasCount+inputCount-1)
	}
	for i := 0; i < biasCount; i++ {
		ng = &NodeGene{Marker: inno.nextMarker(), Type: neural.BIAS, X: step * float64(i), Y: 0}
		genome.Nodes[ng.Marker] = ng
	}

	// Create the input nodes
	for i := 0; i < inputCount; i++ {
		ng = &NodeGene{Marker: inno.nextMarker(), Type: neural.INPUT, X: step * float64(i+biasCount), Y: 0}
		genome.Nodes[ng.Marker] = ng
	}

	// Create the output nodes
	step = 0
	if outputCount > 1 {
		step = 1.0 / float64(outputCount-1)
	}
	for i := 0; i < outputCount; i++ {
		ng = &NodeGene{Marker: inno.nextMarker(), Type: neural.OUTPUT, X: step * float64(i), Y: 1.0}
		genome.Nodes[ng.Marker] = ng
	}

	// Create the connections
	for _, in := range genome.Nodes {
		for _, out := range genome.Nodes {
			if out.Type == neural.OUTPUT && (in.Type == neural.BIAS || in.Type == neural.INPUT) {
				cg := &ConnGene{Marker: inno.nextMarker(),
					Enabled: true, Weight: 0, Source: in.Marker,
					Target: out.Marker}
				genome.Conns[cg.Marker] = cg
			}
		}
	}

	return

}
