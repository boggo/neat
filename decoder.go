package neat

import (
    "github.com/boggo/neural"
    "sort"
)

// Interface describing method to decode a genome into a phenome
type GenomeDecoder interface {
    Decode(genome *Genome) (phenome *Phenome, err error)
}

// Default NEAT decoder
type Decoder struct{}

// Decodes a NEAT genome into a phenome
func (decoder Decoder) Decode(genome *Genome) (phenome Phenome, err error) {

    // Sort the genes to ensure the network is built in the correct order
    n := sortNodes{genome}
    sort.Sort(n)

    c := sortConns{genome}
    sort.Sort(c)

    // Construct the Network.
    network := &neural.Network{}

    // Add the Nodes
    m := make(map[uint32]neural.Node, len(genome.NodeGenes)) // Map of node to gene's marker
    for _, x := range genome.NodeGenes {
        node := neural.NewNode(x.FuncType, x.NodeType)
        m[x.Marker] = node
        network.AddNode(node)
    }

    // Add the Connections. Make sure the connections are sorted so they fire
    // in the correct order
    conns := genome.ConnGenes
    for _, x := range conns {
        if x.Enabled {
            conn := neural.NewConnection(m[x.Source], m[x.Target], x.Weight)
            network.AddConnection(conn)
        }
    }

    // Return the new phenome
    phenome = network
    return

}

// Helper struct to sort Node Genese by position
type sortNodes struct {
    genome *Genome
}

// Len is the number of elements in the collection.
func (s sortNodes) Len() int { return len(s.genome.NodeGenes) }

// Less returns whether the element with index i should sort
// before the element with index j.
func (s sortNodes) Less(i, j int) bool {
    a := s.genome.NodeGenes[i].Position
    b := s.genome.NodeGenes[j].Position

    for k := 0; k < len(a); k++ {
        if a[k] < b[k] {
            return true
        } else if a[k] > b[k] {
            return false
        }
    }
    return false // These two have the same fitness
}

// Swap swaps the elements with indexes i and j.
func (s sortNodes) Swap(i, j int) {
    s.genome.NodeGenes[i], s.genome.NodeGenes[j] = s.genome.NodeGenes[j], s.genome.NodeGenes[i]
}

// Helper struct to sort Connection Genes by their Target Node Gene's
// position.
type sortConns struct {
    genome *Genome
}

// Len is the number of elements in the collection.
func (s sortConns) Len() int { return len(s.genome.ConnGenes) }

// Less returns whether the element with index i should sort
// before the element with index j.
func (s sortConns) Less(i, j int) bool {

    // This is function should be called when the node genes are already sorted.
    // Therefore we can just compare the relative indexes of the two connection
    // genes' target gene nodes
    var a, b int
    for p, x := range s.genome.NodeGenes {
        if x.Marker == s.genome.ConnGenes[i].Target {
            a = p
        }
        if x.Marker == s.genome.ConnGenes[j].Target {
            b = p
        }
    }

    return a < b
}

// Swap swaps the elements with indexes i and j.
func (s sortConns) Swap(i, j int) {
    s.genome.ConnGenes[i], s.genome.ConnGenes[j] = s.genome.ConnGenes[j], s.genome.ConnGenes[i]
}
