package neat

/*

This source file contains the functions and methods for mutating Genomes and their Genes.

*/

import (
    "fmt"
    "github.com/boggo/neural"
    "math/rand"
)

// Mutates a Genome. This function will call the other mutation function based on
// the probabilities set out in the Settings. The random number generator is passed
// in as an argument to improve performance (instead of recreating a new one on each
// call) but this assumes that the mutation occurs in the same go routine as the one
// that created the random number generator
func mutateGenome(gen *Genome, rng *rand.Rand, settings *Settings, inno chan innovator) {
    fmt.Println("Mutate Genome ", gen.ID)

    // Mutate the existing Connection Genes
    if rng.Float64() < settings.MutateWeight {
        for i, _ := range gen.ConnGenes {

            // Mutate the Weight
            if rng.Float64() < settings.MutateWeight {
                mutateWeight(&gen.ConnGenes[i], rng,
                    (rng.Float64() < settings.MutateWeightNew))
            }

            // Mutate Enabled
            if rng.Float64() < settings.MutateEnabled {
                mutateEnabled(&gen.ConnGenes[i])
            }
        }
    }

    // Mutate the existing Node Genes
    for i, _ := range gen.NodeGenes {
        if rng.Float64() < settings.MutateFuncType {
            mutateFunc(&gen.NodeGenes[i], rng)
        }
    }

    // Mutate the structure by adding Genes
    if rng.Float64() < settings.MutateAddNode {
        mutateAddNode(gen, rng, inno)
    }

    if rng.Float64() < settings.MutateAddConnection {
        mutateAddConn(gen, rng, inno)
    }

}

// Mutates the function type Node Gene.
func mutateFunc(node *NodeGene, rng *rand.Rand) {
    f := rng.Intn(len(neural.FuncTypes)-1) + 1 // skip the DIRECT function type
    node.FuncType = neural.FuncTypes[f]
}

// Mutates a ConnGene's Weight. newWeight value determines if the weight
// should be replaced with a whole new value or just perturbed.
func mutateWeight(conn *ConnGene, rng *rand.Rand, newWeight bool) {
    if newWeight {
        conn.Weight = rng.NormFloat64()
    } else {
        conn.Weight = conn.Weight + (rng.Float64() * 2.0) - 1.0
    }
}

// Mutates the Connection Gene's Enabled property by toggling it to the othere
// boolean value
func mutateEnabled(conn *ConnGene) {
    conn.Enabled = !conn.Enabled
}

// Mutates the Genome by adding a new NodeGene. An existing connection is split
// and the new node is inserted in between. The inboound connection is given a
// weight of 1 while the outbound connection preserves the weight of the
// original connection.
func mutateAddNode(gen *Genome, rng *rand.Rand, inno chan innovator) {

    // Identify the connection to split
    i := rng.Int31n(int32(len(gen.ConnGenes)))
    old := gen.ConnGenes[i]
    var src, tgt NodeGene
    for _, x := range gen.NodeGenes {
        if x.Marker == old.Source {
            src = x
        }
        if x.Marker == old.Target {
            tgt = x
        }
    }

    // Borrow the innnovator
    v := <-inno

    // Identify the Position of this new NodeGene. If node already exists within
    // THIS Genome then we simply add a small increment to the Z dimension
    pos := [3]float32{((src.Position[0] + tgt.Position[0]) / 2.0),
        ((src.Position[1] + tgt.Position[1]) / 2.0), 0}
    for _, x := range gen.NodeGenes {
        if x.Position[0] == pos[0] && x.Position[1] == pos[1] {
            if x.Position[2] > pos[2] {
                pos[2] = x.Position[2] + 0.0001
            }
        }
    }

    // Create a new NodeGene
    f := rng.Intn(len(neural.FuncTypes)-1) + 1 // skip the DIRECT function type
    node := NodeGene{NodeType: neural.HIDDEN, FuncType: neural.FuncTypes[f], Position: pos}
    v.blessNodeGene(&node)
    gen.addNodeGene(node)

    // Create the new ConnGenes
    ic := ConnGene{Enabled: old.Enabled, Weight: 1, Source: src.Marker, Target: node.Marker}
    v.blessConnGene(&ic)
    gen.addConnGene(ic)

    oc := ConnGene{Enabled: old.Enabled, Weight: old.Weight, Source: node.Marker, Target: tgt.Marker}
    v.blessConnGene(&oc)
    gen.addConnGene(oc)

    // Return the innovator
    inno <- v

    // Disable the old gene
    old.Enabled = false
}

// Mutates a Genome by adding a new connection between two existing Nodes.
func mutateAddConn(gen *Genome, rng *rand.Rand, inno chan innovator) {

    // Identify the source and target NodeGenes for this new connection
    s := rng.Intn(len(gen.NodeGenes))
    t := rng.Intn(len(gen.NodeGenes))
    if gen.NodeGenes[s].Position[0] > gen.NodeGenes[t].Position[0] {
        s, t = t, s // Make sure the source is below or at the same level as the target
    }

    // Create the connection and added it to the Genome
    w := rng.NormFloat64()
    conn := ConnGene{Enabled: true, Weight: w, Source: gen.NodeGenes[s].Marker,
        Target: gen.NodeGenes[t].Marker}
    v := <-inno
    v.blessConnGene(&conn)
    inno <- v
    gen.addConnGene(conn)
}
