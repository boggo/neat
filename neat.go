package neat

import (
    "math/rand"
    "sort"
    "time"
)

func Load(arch Archiver) (exp *Experiment, err error) {

    exp, err = arch.Restore()
    exp.initSequences()
    return
}

func Go(exp *Experiment, arch Archiver) {

    // Iterate indefinitely
    maxIter := 5 // Can set this to really high number like,
    for i := 0; i < maxIter; i++ {

        // Initialize or advance the generation
        if exp.Current == nil || exp.Current.Species == nil {
            exp.Current = firstPop(exp.Settings)
        } else {
            exp.Current = rollPop(exp.Current, exp.Settings)
        }

        // Decode the genomes
        i := 0
        exp.Current.phenomes = make([]*Phenome, exp.Settings.PopulationSize)
        for s, _ := range exp.Current.Species {
            for g, _ := range exp.Current.Species[s].Genomes {
                p, err := exp.decoder.Decode(g)
                if err != nil {
                    panic(err)
                }
                exp.Current.phenomes[i] = p
            }
        }

        // Evaluate the population
        exp.popEval.Evaluate(exp.Current)

        // Archive the population
        arch.Archive(exp)

    }
}




// TODO: Turn this into a service like sequences but deliver a new Rand each time. Aim is to prevent calls within the same nano second and, thus, the same seed. Probably rare but possible.
func newRandom() *rand.Rand {
    return rand.New(rand.NewSource(time.Now().UnixNano()))
}
