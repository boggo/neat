package neat

import (
    "math/rand"
    "time"
	"fmt"
)

// Run the experiment for a specific number of iterations. May use math.MaxInt32
// or some other really high number to simulate an continuous loop
func Run(exp *Experiment, arch Archiver, iters int) (err error) {

    // Iterate indefinitely
    a := 0
    for t := 0; t < iters; t++ {
		fmt.Println("Run iteration", t)
		
        // Initialize or advance the generation
        if exp.Current == nil || exp.Current.Species == nil {
            exp.Current, err = initPop(exp.Settings)
			if err != nil {
				panic(err)
			}
        } else {
            exp.Current, err = rollPop(exp.Current, exp.Settings)
			if err != nil {
				panic(err)
			}
        }

        // Decode the genomes
        i := 0
        exp.Current.orgs = make([]*Organism, exp.Settings.PopulationSize)
        for s, _ := range exp.Current.Species {
            for g, _ := range exp.Current.Species[s].Genomes {
                p, err := exp.decoder.Decode(exp.Current.Species[s].Genomes[g])
                if err != nil {
                    panic(err)
                }
                exp.Current.orgs[i] = &Organism{p, exp.Current.Species[s].Genomes[g]}
                i++
            }
        }

        // Evaluate the population
        exp.popEval.Evaluate(exp.Current, exp.orgEval)

        // Archive the population
        a++
        if a > exp.Settings.ArchiveFreq {
            arch.Archive(exp)
            a = 0
        }

    }
	
	return
}

// TODO: Turn this into a service like sequences but deliver a new Rand each time. Aim is to prevent calls within the same nano second and, thus, the same seed. Probably rare but possible.
func newRandom() *rand.Rand {
    return rand.New(rand.NewSource(time.Now().UnixNano()))
}
