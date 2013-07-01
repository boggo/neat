package neat

import (
    "github.com/boggo/neural"
    "sort"
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
    Position [3]float32      // 3D position of node in network. In Y,X,Z order for easy sorting for Decoding
}

// A slice of Node Genes.
type NodeGeneSlice []NodeGene

// Len is the number of elements in the collection.
func (s NodeGeneSlice) Len() int { return len(s) }

// Less returns whether the element with index i should sort
// before the element with index j.
func (s NodeGeneSlice) Less(i, j int) bool { return s[i].Marker < s[j].Marker }

// Swap swaps the elements with indexes i and j.
func (s NodeGeneSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

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

// Len is the number of elements in the collection.
func (s ConnGeneSlice) Len() int { return len(s) }

// Less returns whether the element with index i should sort
// before the element with index j.
func (s ConnGeneSlice) Less(i, j int) bool { return s[i].Marker < s[j].Marker }

// Swap swaps the elements with indexes i and j.
func (s ConnGeneSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

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

// Determine the maximum Marker within a genome
func (g *Genome) maxMarker() uint32 {
    n := g.NodeGenes[len(g.NodeGenes)-1].Marker
    c := g.ConnGenes[len(g.ConnGenes)-1].Marker
    if n > c {
        return n
    } else {
        return c
    }
}

func (gen *Genome) addNodeGene(node NodeGene) {
    gen.NodeGenes = append(gen.NodeGenes, node)
    sort.Sort(gen.NodeGenes)
}

func (gen *Genome) addConnGene(conn ConnGene) {
    gen.ConnGenes = append(gen.ConnGenes, conn)
    sort.Sort(gen.ConnGenes)
}
