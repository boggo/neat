package neat

import (
    "github.com/boggo/neural"
	"github.com/boggo/sequence"
	"math"
)

// Genes are individual elements within the Genome. Each are tagged with a
// historical Marker indicating when they became part of the Genome.
type Gene struct {
    Marker uint32
}

// Node Genes encode the information necessary for constucting Neural Network
// nodes (neurons).
type NodeGene struct {
    Gene                     // Common gene components like historical marker
    NodeType neural.NodeType // Type of node
    FuncType neural.FuncType // Activation function
    Position [3]float32      // 3D position of node in network. In Y,X,Z order for easy sorting
}

// A slice of Node Genes.
type NodeGeneSlice []NodeGene

// Connection Genes encode the information for constructing Neural Network
// connections (synapses)
type ConnGene struct {
    Gene            // Common gene components like historical marker
    Enabled bool    // Is this connection enabled
    Weight  float64 // The weight applied during activcation
    Source  uint32  // Historical Marker of the source node
    Target  uint32  // Historical Marker of the target node
}

// A slice of Connection Genes.
type ConnGeneSlice []ConnGene

// Genome encodes the information necessary to construct a Network for
// use in an Experiment
type Genome struct {
    ID        uint32        // Uniquely identifies this Genome
    NodeGenes NodeGeneSlice // Node definitions
    ConnGenes ConnGeneSlice // Connection definitions
    Fitness   []float64     // Fitness of this Genome
}

// A slice of Genomes
type GenomeSlice []*Genome

// Len is the number of elements in the collection.
func (g GenomeSlice) Len() int { return len(g) }

// Less returns whether the element with index i should sort
// before the element with index j.
func (g GenomeSlice) Less(i, j int) bool {
    // Order Genome by their relative fitness
    return g[i].Fitness[0] < g[j].Fitness[0]
}

// Swap swaps the elements with indexes i and j.
func (g GenomeSlice) Swap(i, j int) { g[i], g[j] = g[j], g[i] }


// Calcualtes the compatiblity distance between two genomes. c1, c2 and
// c3 refer to the coefficients described in section 3.3 of "Evolving
// Neural Networks through Augmented Topologies" (Stanley et. al., P109)
func distance(genome1, genome2 *Genome, settings *Settings) float64 {

    // The working variables
    var e, d, n, w, s float64

    // Identify the maximum marker in second genome to determine
    // which genes are disjoint and which are excess
    m := maxMarker(genome2)

    // Process the node genes
    for _, g1 := range genome1.NodeGenes {
        found := false
        for _, g2 := range genome2.NodeGenes {
            if g1.Marker == g2.Marker {
                found = true
                break
            }
        }
        if !found {
            if g1.Marker < m {
                d++ // disjoint
            } else {
                e++ // excess
            }
        }
    }

    // Process the connection genes
    for _, g1 := range genome1.ConnGenes {
        found := false
        for _, g2 := range genome2.ConnGenes {
            if g1.Marker == g2.Marker {
                found = true
                w += math.Abs(g1.Weight - g2.Weight)
                s += 1.0
                break
            }
        }
        if !found {
            if g1.Marker < m {
                d++ // disjoint
            } else {
                e++ // excess
            }
        }
    }

    // Calculate the distance between the two genomes
    c1 := settings.ExcessCoefficient
    c2 := settings.DisjointCoefficient
    c3 := settings.WeightCoefficient
    a := len(genome1.NodeGenes) + len(genome1.ConnGenes)
    b := len(genome2.NodeGenes) + len(genome2.ConnGenes)
    n = math.Max(float64(a), float64(b))
    w = w / s
    return c1*e/n + c2*d/n + c3*w
}

// Determine the maximum Marker within a genome
func maxMarker(g *Genome) uint32 {
    m := uint32(0)
    for _, x := range g.NodeGenes {
        if m < x.Marker {
            m = x.Marker
        }
    }
    for _, x := range g.ConnGenes {
        if m < x.Marker {
            m = x.Marker
        }
    }
    return m
}

// Add a node gene
func addNodeGene(genome *Genome, node NodeGene) {
    genome.NodeGenes = append(genome.NodeGenes, node)
}

// Add a connection gene
func addConnGene(genome *Genome, conn ConnGene) {
    genome.ConnGenes = append(genome.ConnGenes, conn)
}

func initialGenome(settings Settings, marker *sequence.UInt32) (genome *Genome) {
	
    // Shortcuts to settings values
    biasCount := settings.BiasCount
    inputCount := settings.InputCount
    outputCount := settings.OutputCount

    // Note the offsets
    nodeCount := biasCount + inputCount + outputCount
    inputOffset := biasCount
    outputOffset := inputOffset + inputCount

    // Create a new Genome
    genome = &Genome{NodeGenes: make([]NodeGene, nodeCount),
        ConnGenes: make([]ConnGene, (biasCount+inputCount)*outputCount)}

    // Create the bias and input nodes
    var step float32
    step = 0
    if biasCount+inputCount > 1 {
        step = 1.0 / float32(biasCount+inputCount-1)
    }
    for i := 0; i < biasCount; i++ {
        pos := [3]float32{0, step * float32(i), 0}
        genome.NodeGenes[i] = NodeGene{Gene: Gene{Marker: marker.Next()},
            NodeType: neural.BIAS, FuncType: neural.DIRECT, Position: pos}
    }

    // Create the input nodes
    for i := 0; i < inputCount; i++ {
        pos := [3]float32{0, float32(i+inputOffset) * step, 0}
        genome.NodeGenes[i+inputOffset] = NodeGene{Gene: Gene{Marker: marker.Next()},
            NodeType: neural.INPUT, FuncType: neural.DIRECT, Position: pos}
    }

    // Create the output nodes
    step = 0
    if outputCount > 1 {
        step = 1.0 / float32(outputCount-1)
    }
    for i := 0; i < outputCount; i++ {
        pos := [3]float32{1, float32(i) * step, 0}
        genome.NodeGenes[i+outputOffset] = NodeGene{Gene: Gene{Marker: marker.Next()},
            NodeType: neural.OUTPUT, FuncType: neural.STEEPENED_SIGMOID, Position: pos}
    }

    // Create the connections
    for i := 0; i < biasCount+inputCount; i++ {
        for o := 0; o < outputCount; o++ {
            genome.ConnGenes[i+o] = ConnGene{Gene: Gene{Marker: marker.Next()},
                Enabled: true, Weight: 0, Source: genome.NodeGenes[i].Marker,
                Target: genome.NodeGenes[o+outputOffset].Marker}
        }
    }

    // Return the genome
    return
	
}