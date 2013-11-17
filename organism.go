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
	"github.com/boggo/neural"
	"math"
)

type Phenome interface {
	Analyze(inputs []float64) (outputs []float64, err error)
}

type Organism struct {
	*Genome
	Phenome `json:"-"`
}

func cloneOrg(source *Organism, id int) (clone *Organism) {
	clone = &Organism{Genome: cloneGenome(source.Genome, id)}
	// phenome will be decoded during next iteration
	return
}

func mutate(settings *Settings, inno *innovation, org *Organism) {

	switch {
	case random.Next() < settings.MutateAddNode:
		mutateAddNode(inno, org)
	case random.Next() < settings.MutateAddConnection:
		mutateAddConn(settings, inno, org)
	default:
		for _, cg := range org.Conns {
			if random.Next() < settings.MutateWeight {
				if random.Next() < settings.MutateWeightNew {
					mutateWeightNew(cg)
				} else {
					mutateWeight(cg)
				}
			}
			if random.Next() < settings.MutateEnabled {
				mutateEnabled(cg)
			}
		}
	}
}

func mutateAddNode(inno *innovation, org *Organism) {

	// Pick a connection to split
	var old *ConnGene
	i := random.Int(len(org.Conns))
	j := 0
	for _, v := range org.Conns {
		if i == j {
			old = v
			break
		}
		j += 1
	}

	// Note the old source and target
	src := org.Nodes[old.Source]
	tgt := org.Nodes[old.Target]

	// Create a new node
	ng := &NodeGene{Type: neural.HIDDEN, X: (src.X + tgt.X) / 2.0, Y: (src.Y + tgt.Y) / 2.0}
	ng.Marker = inno.blessNodeGene(nodeKey{ng.X, ng.Y})
	org.Nodes[ng.Marker] = ng

	// Create the new connections
	cg1 := &ConnGene{Source: src.Marker, Target: ng.Marker, Enabled: true, Weight: 1.0}
	cg1.Marker = inno.blessConnGene(connKey{cg1.Source, cg1.Target})
	org.Conns[cg1.Marker] = cg1
	cg2 := &ConnGene{Source: ng.Marker, Target: tgt.Marker, Enabled: true, Weight: old.Weight}
	cg2.Marker = inno.blessConnGene(connKey{cg2.Source, cg2.Target})
	org.Conns[cg2.Marker] = cg2

	// Disable the old connection
	old.Enabled = false
}

func mutateAddConn(settings *Settings, inno *innovation, org *Organism) {

	// Pick 2 nodes to connect
	var ng1, ng2 *NodeGene
	a := random.Int(len(org.Nodes))
	b := random.Int(len(org.Nodes)-settings.BiasCount-settings.InputCount) +
		settings.BiasCount + settings.InputCount
	j := 0
	for _, v := range org.Nodes {
		if a == j {
			ng1 = v
		}
		if b == j {
			ng2 = v
		}
		j++
	}

	// validate the nodes
	if ng1.Marker == ng2.Marker {
		return // No connections to the same node
	}
	if ng1.Y > ng2.Y {
		ng1, ng2 = ng2, ng1
	}
	if ng1.Type == neural.OUTPUT {
		return
	}
	if ng2.Type == neural.BIAS || ng2.Type == neural.INPUT {
		return
	}

	// Look for an existing connection between these nodes
	found := false
	for _, c := range org.Conns {
		if c.Source == ng1.Marker && c.Target == ng2.Marker {
			found = true
			break
		}
	}
	if found {
		return // we already have this connection
	}

	// Make the new connection
	cg := &ConnGene{Source: ng1.Marker, Target: ng2.Marker, Enabled: true, Weight: random.Gaussian()}
	cg.Marker = inno.blessConnGene(connKey{cg.Source, cg.Target})
	org.Conns[cg.Marker] = cg
}

func mutateWeight(cg *ConnGene) {
	cg.Weight += random.Gaussian()
	if cg.Weight > 30.0 {
		cg.Weight = 30
	}
	if cg.Weight < -30.0 {
		cg.Weight = -30
	}
}

func mutateWeightNew(cg *ConnGene) {
	cg.Weight = random.Gaussian()
}

func mutateEnabled(cg *ConnGene) {
	cg.Enabled = true
}

func crossover(inno *innovation, p1, p2 *Organism) (child *Organism) {

	// Order parents by fitness
	if p2.Fitness[0] > p1.Fitness[0] {
		p1, p2 = p2, p1
	}

	// Create the new child
	genome := &Genome{ID: inno.nextID(), Nodes: make(map[int]*NodeGene), Conns: make(map[int]*ConnGene)}
	child = &Organism{Genome: genome}

	// Crossover the connection genes
	for _, cg1 := range p1.Conns {
		cg2, ok := p2.Conns[cg1.Marker]
		if ok {
			if random.Next() < 0.5 {
				child.Conns[cg1.Marker] = cloneConn(cg1)
			} else {
				child.Conns[cg2.Marker] = cloneConn(cg2)
			}
		} else {
			child.Conns[cg1.Marker] = cloneConn(cg1)
		}
	}

	// Crossover the node genes
	var ng1, ng2 *NodeGene
	var ok bool
	for _, cg1 := range child.Conns {
		_, ok = child.Nodes[cg1.Source] // look first in child
		if !ok {
			ng1, ok = p1.Nodes[cg1.Source] // Grab from parent 1
			if ok {
				ng2, ok = p2.Nodes[cg1.Source] // Grab from parent 2
				if ok {
					if random.Next() < 0.5 {
						child.Nodes[ng1.Marker] = cloneNode(ng1)
					} else {
						child.Nodes[ng2.Marker] = cloneNode(ng2)
					}
				} else {
					child.Nodes[ng1.Marker] = cloneNode(ng1)
				}
			} else {
				child.Nodes[ng2.Marker] = cloneNode(ng2)
			}
		}
		_, ok = child.Nodes[cg1.Target] // look first in child
		if !ok {
			ng1, ok = p1.Nodes[cg1.Target] // Grab from parent 1
			if ok {
				ng2, ok = p2.Nodes[cg1.Target] // Grab from parent 2
				if ok {
					if random.Next() < 0.5 {
						child.Nodes[ng1.Marker] = cloneNode(ng1)
					} else {
						child.Nodes[ng2.Marker] = cloneNode(ng2)
					}
				} else {
					child.Nodes[ng1.Marker] = cloneNode(ng1)
				}
			} else {
				child.Nodes[ng2.Marker] = cloneNode(ng2)
			}
		}
	}
	return
}

// Returns the compatibiliy distance between the two organisms
func distance(settings *Settings, o1, o2 *Organism) float64 {

	// Note from http://www.cs.ucf.edu/~kstanley/neat.html
	// (How should I test my own version of NEAT to make sure it works)
	//
	// If you decide to use the species compatibility coefficients and thresholds
	// from my own .ne settings files (provided with my NEAT release), then do not
	// normalize the terms in the compatibility function, because I did not do this
	// with my .ne files. In other words, even though my papers suggest normalizing
	// (dividing my number of genes), since I didn't do that the coefficients that
	// I used will not work the same for you if you normalize. If you strongly desire
	// to normalize, you will need to find your own appropriate coefficients and
	// threshold.

	// To use the default settings from Stanley's paper we only consider conn genes.
	// Look first at the organism with the most conn genes
	if len(o1.Conns) < len(o2.Conns) {
		o1, o2 = o2, o1
	}
	mm := 0 // Max marker in o2
	for _, cg2 := range o2.Conns {
		if cg2.Marker > mm {
			mm = cg2.Marker
		}
	}

	// Make the comparison
	var d, e, m, w float64
	for _, cg1 := range o1.Conns {
		cg2, ok := o2.Conns[cg1.Marker]
		if ok {
			m += 1 // This is a match
			w += math.Abs(cg1.Weight - cg2.Weight)
		} else {
			if cg1.Marker > mm {
				e += 1 // Excess
			} else {
				d += 1 // disjoint
			}
		}
	}
	d += float64(len(o2.Conns)) - m // Look for disjoints in o2

	if m > 0 { // take the average weight difference
		w = w / m
	}

	return settings.ExcessCoefficient*e + settings.DisjointCoefficient*d +
		settings.WeightCoefficient*w
}

type OrganismSlice []*Organism

func (os OrganismSlice) Len() int           { return len(os) }
func (os OrganismSlice) Swap(i, j int)      { os[i], os[j] = os[j], os[i] }
func (os OrganismSlice) Less(i, j int) bool { return os[i].Fitness[0] < os[j].Fitness[0] }

func (os OrganismSlice) TotalFitness() float64 {
	sum := float64(0)
	for _, o := range os {
		sum += o.Fitness[0]
	}
	return sum
}

// Removes a node gene
// From http://sharpneat.sourceforge.net/phasedsearch.html
// Neuron deletion is slightly more complex. The deletion algorithm attempts to replace neurons with connections to maintain any circuits a neuron may have participated in, in further generations those connections themselves will be open to deletion. This approach provides NEAT with the ability to delete whole structures, not just connections.
//
//Because we replace connected neurons with connections we must be careful which neurons we delete. Any neuron with only incoming or only outgoing connections is at a dead-end of a circuit and can therefore be safely deleted with all of it's connections. However, a neuron with multiple incoming and multiple outgoing connections will require a large number of connections to substitute for the loss of the neuron - we must fully connect all of the original neuron's source neurons with its target neurons, this could be done but may actually be detrimental since the functionality represented by the neuron is now distributed over a number of connections, and this cannot easily be reversed. Because of this, such neurons are omitted from the process of selecting neurons for deletion.
//
// Neurons with only one incoming or one outgoing connection can be replaced with however many connections were on the other side of the neuron, therefore these are candidates for deletion.

func mutateDelNode(settings *Settings, org *Organism) {

	// Pick a node to delete
	var n *NodeGene
	i := random.Int(len(org.Nodes))
	for _, v := range org.Nodes {
		n = v
		if i == 0 {
			break
		}
		i -= 1
	}
	if n.Type != neural.HIDDEN {
		return
	} // Only remove hidden nodes

	// Node the incoming and outgoing connections
	incoming := make([]*ConnGene, 0, 10)
	outgoing := make([]*ConnGene, 0, 10)
	for _, c := range org.Conns {
		if c.Source == n.Marker {
			outgoing = append(outgoing, c)
		}
		if c.Target == n.Marker {
			incoming = append(incoming, c)
		}
	}

	// The node is cut off (no incoming or no outgoing connections)
	if len(incoming) == 0 {
		for _, c := range outgoing {
			delete(org.Conns, c.Marker)
		}
		delete(org.Nodes, n.Marker)
	} else if len(outgoing) == 0 {
		for _, c := range incoming {
			delete(org.Conns, c.Marker)
		}
		delete(org.Nodes, n.Marker)
	}

	// There is only one incoming node
	if len(incoming) == 1 {

		// Replace the node in the outgoing connections
		a := incoming[0]
		for _, c := range outgoing {
			c.Source = a.Source
		}

		// Delete the incoming connection and node
		delete(org.Conns, a.Marker)
		delete(org.Nodes, n.Marker)

	} else if len(outgoing) == 1 {

		// Replace the node in the incoming connections
		a := outgoing[0]
		for _, c := range incoming {
			c.Target = a.Target
		}

		// Delete the incoming connection and node
		delete(org.Conns, a.Marker)
		delete(org.Nodes, n.Marker)
	}
}

// Removes a connection gene
// From http://sharpneat.sourceforge.net/phasedsearch.html
// Connection deletion is very simply the deletion of a randomly selected connection, all connections are considered to be available for deletion. When a connection is deleted the neurons that were at each end of the connection are tested to check if they are no longer connected to by other connections, if this is the case then the stranded neuron is also deleted. Note that a more thorough cleanup routine could be invoked at this point that cleans up any dead-end structures that could not possibly be functional, but this can become complex and so we leave NEAT to eliminate such structures naturally.
//
func mutateDelConnection(settings *Settings, org *Organism) {

	// Pick a connection to remove
	if len(org.Conns) == 0 {
		return
	}
	var c *ConnGene
	i := random.Int(len(org.Conns))
	for _, v := range org.Conns {
		c = v
		if i == 0 {
			break
		}
		i -= 1
	}

	// Node the nodes connected
	src := org.Nodes[c.Source]
	tgt := org.Nodes[c.Target]
	sok := src.Type == neural.HIDDEN // source ok to delete as well
	tok := tgt.Type == neural.HIDDEN // target ok to delete as well
	for k, v := range org.Conns {
		if k != c.Marker {
			sok = sok && !(v.Source == src.Marker || v.Target == src.Marker)
			tok = tok && !(v.Source == tgt.Marker || v.Target == tgt.Marker)
		}
	}

	// Remove the nodes
	if sok {
		delete(org.Nodes, src.Marker)
	}
	if tok {
		delete(org.Nodes, tgt.Marker)
	}

	// Remove the connection
	delete(org.Conns, c.Marker)
}
