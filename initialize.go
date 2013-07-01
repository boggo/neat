package neat

/*

This source file contains the functions and methods for creating the initial
structs and environment

*/

import (
    "fmt"
    "github.com/boggo/neural"
    "github.com/boggo/sequence"
)

// Creates the first Population, creating the prototypical Genome and
// then cloning it to fill up the first Species with PopulationSize copies.
func initPop(settings *Settings) (pop *Population, err error) {
    fmt.Println("Initial Population")

    // Create the initial Species
    spec, err := initSpec(settings)
    if err != nil {
        panic(err)
    }

    // Create the new generation
    popSize := settings.PopulationSize
    pop = &Population{Generation: 0,
        Species: []*Species{spec},
        orgs:    make([]*Organism, popSize)}

    // Create the prototype genome
    gen, err := initGen(settings)
    settings.ProtoGenome = *gen // Record the Genome for future use

    // Add genomes to the species and the phenomes to the population
    r := newRandom()
    for i := 0; i < popSize; i++ {

        // Create a new genome
        g := cloneGenome(gen, settings.ids.Next())

        // Randomize the weights
        for w, _ := range g.ConnGenes {
            g.ConnGenes[w].Weight = r.Float64()*2.0 - 1.0
        }

        // Add to the species
        pop.Species[0].Genomes[i] = g

    }

    return
}

// Creates the first Species to hold the entire first PopulationSize of
// Genomes
func initSpec(settings *Settings) (spec *Species, err error) {
    fmt.Println("Initial Species")

    spec = &Species{Genomes: make([]*Genome, settings.PopulationSize)}
    return
}

// Creates the initial Genome by fully connecting the bias and input nodes
// to the output node(s)
func initGen(settings *Settings) (gen *Genome, err error) {
    fmt.Println("Initial Genome")

    // Shortcuts to settings values
    biasCount := settings.BiasCount
    inputCount := settings.InputCount
    outputCount := settings.OutputCount

    // Note the offsets
    nodeCount := biasCount + inputCount + outputCount
    inputOffset := biasCount
    outputOffset := inputOffset + inputCount

    // Create a new Genome
    gen = &Genome{NodeGenes: make([]NodeGene, nodeCount),
        ConnGenes: make([]ConnGene, (biasCount+inputCount)*outputCount)}

    // Create the bias and input nodes
    var step float32
    step = 0
    if biasCount+inputCount > 1 {
        step = 1.0 / float32(biasCount+inputCount-1)
    }
    for i := 0; i < biasCount; i++ {
        pos := [3]float32{0, step * float32(i), 0}
        gen.NodeGenes[i] = NodeGene{Gene: Gene{Marker: settings.marker.Next()},
            NodeType: neural.BIAS, FuncType: neural.DIRECT, Position: pos}
    }

    // Create the input nodes
    for i := 0; i < inputCount; i++ {
        pos := [3]float32{0, float32(i+inputOffset) * step, 0}
        gen.NodeGenes[i+inputOffset] = NodeGene{Gene: Gene{Marker: settings.marker.Next()},
            NodeType: neural.INPUT, FuncType: neural.DIRECT, Position: pos}
    }

    // Create the output nodes
    step = 0
    if outputCount > 1 {
        step = 1.0 / float32(outputCount-1)
    }
    for i := 0; i < outputCount; i++ {
        pos := [3]float32{1, float32(i) * step, 0}
        gen.NodeGenes[i+outputOffset] = NodeGene{Gene: Gene{Marker: settings.marker.Next()},
            NodeType: neural.OUTPUT, FuncType: neural.STEEPENED_SIGMOID, Position: pos}
    }

    // Create the connections
    for i := 0; i < biasCount+inputCount; i++ {
        for o := 0; o < outputCount; o++ {
            gen.ConnGenes[i+o] = ConnGene{Gene: Gene{Marker: settings.marker.Next()},
                Enabled: true, Weight: 0, Source: gen.NodeGenes[i].Marker,
                Target: gen.NodeGenes[o+outputOffset].Marker}
        }
    }

    // Return the genome
    return
}

// Itialize the sequences based on the last known values. This is run after
// restoring an experiment from an archive
func initSeq(exp *Experiment) {
    fmt.Println("Initial Sequences")

    // Iterate the species in the current generation
    i := uint32(0)
    m := uint32(0)
    if exp.Current != nil && exp.Current.Species != nil {
        for _, s := range exp.Current.Species {

            // Iterate the genomes in the species
            for _, g := range s.Genomes {

                // Increment the last know ID
                if g.ID > i {
                    i = g.ID
                }

                // Identify the last marker
                mm := g.maxMarker()
                if mm > m {
                    m = mm
                }

            }
        }
        // Increment to begin at next value
        i += 1
        m += 1

    }

    // Start the services
    exp.Settings.ids = sequence.NewUInt32(i, 1)
    exp.Settings.marker = sequence.NewUInt32(m, 1)

}

func cloneGenome(orig *Genome, id uint32) (clone *Genome) {
    clone = &(*orig) // take a copy of the original
    clone.ID = id    // set the new ID
    return
}
