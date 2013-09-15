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

package decoder

import (
	"github.com/boggo/neat"
	"github.com/boggo/neat/phenome"
	"github.com/boggo/neural" // TODO: Should this library be moved to code.google.com, too?
	"sort"
)

// Default NEAT decoder
type neatDecoder struct{}

// Returns a new NEAT decoder
func NewNEAT() (decoder neat.Decoder) {
	return &neatDecoder{}
}

// Decodes a genome into a phenome using the NEAT decoder
func (d neatDecoder) Decode(genome *neat.Genome) (pnome neat.Phenome, err error) {

	// Extract the nodes into the correct sorted order (by its position)
	nodes := make([]*neat.NodeGene, len(genome.Nodes))
	i := 0
	for _, v := range genome.Nodes {
		nodes[i] = v
		i += 1
	}
	sn := &sortNodes{nodes}
	sort.Sort(sn)

	// Extract the connections into the corrected sorted order (by target node's position)
	conns := make([]*neat.ConnGene, len(genome.Conns))
	i = 0
	for _, v := range genome.Conns {
		conns[i] = v
		i += 1
	}
	sc := &sortConns{genome, conns}
	sort.Sort(sc)

	// Build the network
	network := &neural.Network{}
	nmap := make(map[int]neural.Node)
	for _, ng := range nodes {
		var node neural.Node
		if ng.Type == neural.BIAS || ng.Type == neural.INPUT {
			node = neural.NewNode(neural.DIRECT, ng.Type)
		} else {
			node = neural.NewNode(neural.SIGMOID, ng.Type)
		}
		nmap[ng.Marker] = node
		network.AddNode(node)
	}

	for _, cg := range conns {
		if cg.Enabled {
			var conn neural.Connection
			conn = neural.NewConnection(nmap[cg.Source], nmap[cg.Target], cg.Weight)
			network.AddConnection(conn)
		}
	}

	// Return the phenome
	//network.Dump()
	pnome = phenome.NewNetwork(network)
	return
}

type sortNodes struct {
	nodes []*neat.NodeGene
}

func (s *sortNodes) Len() int { return len(s.nodes) }
func (s *sortNodes) Less(i, j int) bool {
	if s.nodes[i].Y == s.nodes[j].Y {
		return s.nodes[i].X < s.nodes[j].X
	} else {
		return s.nodes[i].Y < s.nodes[j].Y
	}
}
func (s *sortNodes) Swap(i, j int) { s.nodes[i], s.nodes[j] = s.nodes[j], s.nodes[i] }

type sortConns struct {
	genome *neat.Genome
	conns  []*neat.ConnGene
}

func (s *sortConns) Len() int { return len(s.conns) }
func (s *sortConns) Less(i, j int) bool {
	a := s.genome.Nodes[s.conns[i].Target]
	b := s.genome.Nodes[s.conns[j].Target]
	if a.Y == b.Y {
		return a.X < b.X
	} else {
		return a.Y < b.Y
	}
}
func (s *sortConns) Swap(i, j int) { s.conns[i], s.conns[j] = s.conns[j], s.conns[i] }
