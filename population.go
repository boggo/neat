package neat

import (
    "sort"
)

type Population struct {

    // Generation number
    Generation uint32

    // Slice of species - this is what we track
    Species SpeciesSlice

    // Slice of organism - this is what we evaluate
    orgs OrganismSlice
}

func rollPopulation(settings Settings, curr *Population) (next *Population) {

    // Create the next population
    next = newPopulation(curr.Generation+1, len(curr.Species),
        settings.PopulationSize)

    // Update each species
    totSpecFit := float64(0)
    for i, _ := range curr.Species {

        // grab a reference to current species and roll it forward
        s := curr.Species[i]
        n := rollSpecies(s, settings)
        next.Species = append(next.Species, n)

        // Update the species's fitness
        totSpecFit += s.fitness

    }

    // Ensure everything is in order
    sort.Sort(curr.Species)
    sort.Reverse(curr.Species)

    sort.Sort(curr.organisms)
    sort.Reverse(curr.organisms)

    // Create offspring
    // TODO: This could be made concurrent if we collect the offspring
    //       in a channel and then speciate afterwards
    created := 0
    allGenomes, totFit := rankGenomes(curr.organisms)
    for i, _ := range curr.Species {

        // Copy the species
        s := curr.Species[i]

        // Determine how many to create
        toCreate := int(s.fitness / totSpecFit * float64(settings.PopulationSize))

        // Peel off any extraneous elite
        if len(next.Species[i].Genomes) > toCreate {
            next.Species[i].Genomes = next.Species[i].Genomes[:toCreate-1]
        }

        // Update created and toCreate
        created += len(next.Species[i].Genomes)
        toCreate -= len(next.Species[i].Genomes)

        // Create the new genomes
        for j := 0; j < toCreate; j++ {
            offspring := createOffspring(s.Genomes, allGenomes, settings, ids, markers)
            // Do we need a chek here for a viable baby? That is, can this genome create a functioning phenome (problem may occur in ESHyperNeat)
            addGenomeToPop(offspring, next, settings) // Adds the genome to the correct species in the new population
            created++
        }

    }

    // Fill the population
    for i := created; i < settings.PopulationSize; i++ {
        offspring := createOffspring(allGenomes, allGenomes, settings)
        // Do we need a chek here for a viable baby? That is, can this genome create a functioning phenome (problem may occur in ESHyperNeat)
        addGenomeToPop(offspring, next, settings) // Adds the genome to the correct species in the new population

    }

    return

}

func rankGenomes(orgs OrganismSlice) (genomes GenomeSlice, totalFit float64) {

    // Create the new slice
    genomes = make([]*Genome, len(orgs))
    for i := 0; i < len(orgs); i++ {
        genomes[i] = orgs[i].Genome
        totalFit += genomes[i].Fitness[0]
    }

    // Reverse sort the slice
    sort.Sort(genomes)
    sort.Reverse(genomes)

    return

}

type PopulationEvaluator interface {
    Evaluate(pop *Population, orgEval OrganismEvaluator) (err error)
}

func NewPopulationEvaluator() PopulationEvaluator {
    return &serialPopEval{}
}

// Default Population Evaluators. Processes organisms one after the other.
// This should be replaced with a concurrent (or even more advanced) evaluator
// for production
type serialPopEval struct{}

// Evaluates the entire population
func (pe serialPopEval) Evaluate(pop *Population, orgEval OrganismEvaluator) (err error) {

    // Iterate the organisms in the population
    // TODO: make a Go routine out of this so it can execute concurrently
    for i, _ := range pop.orgs {
        orgEval.Evaluate(pop.orgs[i])
    }

    // Evaluation complete
    return
}
