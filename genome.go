package neat

import (
    "github.com/boggo/neural"
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
