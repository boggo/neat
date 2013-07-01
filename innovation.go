package neat

import (
    "github.com/boggo/sequence"
)

// Innvovator is an internal structure which is passed through a channel when
// mutating Genomes. This allows a centralized way of adding historical markers
// to new Node and Connection Genes. The fields are set from the Experiement's
// Settings element.
type innovator struct {
    prototype *Genome
    marker    *sequence.UInt32
}

// Blesses a Connection Gene. If the ConnGene has already been added to the
// prototype Genome then the new Gene is given the existing Marker. If not,
// the new Gene is added to the prototype and the new Marker is used.
func (inno *innovator) blessConnGene(conn *ConnGene) {

    // A ConnGene is the same if both the source and target markers match
    // the new Gene
    for _, x := range inno.prototype.ConnGenes {
        if x.Source == conn.Source && x.Target == conn.Target {
            conn.Marker = x.Marker
            return
        }
    }

    // Connection not found, add to the prototype
    conn.Marker = inno.marker.Next()
    inno.prototype.ConnGenes = append(inno.prototype.ConnGenes, *conn) // Store a copy in the prototype

}

// Blesses a Node Gene. If the NodeGene has already been added to the
// prototype Genome then the new Gene is given the existing Marker. If not,
// the new Gene is added to the prototype and the new Marker is used.
func (inno *innovator) blessNodeGene(node *NodeGene) {

    // A NodeGene is the same if its Position a matches the new Gene.
    for _, x := range inno.prototype.NodeGenes {

        // Look for a NodeGene in this Position
		found := true
        for p := 0; p < 3; p++ {
            if x.Position[p] != node.Position[p] {
                found = false
                break
            }
        }

        if found {
            node.Marker = x.Marker
            return
        }
    }

    // Connection not found, add to the prototype
    node.Marker = inno.marker.Next()
    inno.prototype.NodeGenes = append(inno.prototype.NodeGenes, *node) // store a copy in the prototype

}
