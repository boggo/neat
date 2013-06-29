package neat

import (
    "github.com/boggo/sequence"
)

type innovator struct {
    prototype *Genome
    marker    sequence.UInt32
}

func (inno *innovator) blessConnGene(conn *ConnGene) {

    // Look for the same connection
    for _, x := range inno.prototype.ConnGenes {
        if x.Source == conn.Source && x.Target == conn.Target {
            conn.Marker = x.Marker
            return
        }
    }

    // Connection not found, add to the prototype
    conn.Marker = inno.marker.Next()
    addConnGene(inno.prototype, *conn)	// Store a copy in the prototype

}

func (inno *innovator) blessNodeGene(node *NodeGene) {

    // Look for the same node
    for _, x := range inno.prototype.NodeGenes {
        found := x.FuncType == node.FuncType
        if found {
            for p := 0; p < 3; p++ {
                if x.Position[p] != node.Position[p] {
                    found = false
                    break
                }
            }
        }
        if found {
            node.Marker = x.Marker
            return
        }
    }

    // Connection not found, add to the prototype
    node.Marker = inno.marker.Next()
    addNodeGene(inno.prototype, *node)	// store a copy in the prototype

}
