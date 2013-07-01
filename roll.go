package neat

/*

This source file contains the functions and methods for rolling the current
generation into the next one

*/
import (
    "fmt"
    "math"
    "math/rand"
    "sort"
)

// Rolls the Population from one generation to the next
func rollPop(curr *Population, settings *Settings) (next *Population, err error) {
    fmt.Println("rollPop", curr.Generation)
    // Create the next population
    next = &Population{Generation: curr.Generation + 1,
        Species: make([]*Species, len(curr.Species)),
        orgs:    make([]*Organism, settings.PopulationSize)}

    // Update each species
    totSpecFit := float64(0)
    for i, _ := range curr.Species {

        // grab a reference to current species and roll it forward
        rng := newRandom()
        s := curr.Species[i]
        n, err := rollSpecies(s, settings, rng)
        if err != nil {
            panic(err)
        }
        next.Species = append(next.Species, n)

        // Update the species's fitness
        totSpecFit += s.fitness

    }

    // Ensure everything is in order
    sort.Sort(curr.Species)
    sort.Reverse(curr.Species)

    sort.Sort(curr.orgs)
    sort.Reverse(curr.orgs)

    // Create the innovator
    v := innovator{&settings.ProtoGenome, settings.marker}
    inno := make(chan innovator)
    innovate := true
    go func() {
        inno <- v
        for innovate {
        }
    }()

    // Create offspring
    // TODO: This could be made concurrent if we collect the offspring
    //       in a channel and then speciate afterwards
    created := 0
    allGenomes, totFit := rankGenomes(curr.orgs)
    for i, _ := range curr.Species {

        // Copy the species
        s := curr.Species[i]

        // Determine how many to create
        toCreate := int(math.Floor(s.fitness / totSpecFit * float64(settings.PopulationSize)))

        // Peel off any extraneous elite
        // TODO: Rework this bit to make sure there are no memory leaks
        if len(next.Species[i].Genomes) > toCreate {
            next.Species[i].Genomes = next.Species[i].Genomes[:toCreate-1]
        }

        // Update created and toCreate
        created += len(next.Species[i].Genomes)
        toCreate -= len(next.Species[i].Genomes)

        // Create the new genomes
        rng := newRandom()
        for j := 0; j < toCreate; j++ {
            offspring := createOffspring(s.Genomes, s.fitness, allGenomes, totFit, settings, rng, inno)
            // Do we need a chek here for a viable baby? That is, can this genome create a functioning phenome (problem may occur in ESHyperNeat)
            addGenomeToPop(offspring, next, settings) // Adds the genome to the correct species in the new population
            created++
        }

    }

    // Fill the population
    rng2 := newRandom()
    for i := created; i < settings.PopulationSize; i++ {
        offspring := createOffspring(allGenomes, totFit, allGenomes, totFit, settings, rng2, inno)
        // Do we need a chek here for a viable baby? That is, can this genome create a functioning phenome (problem may occur in ESHyperNeat)
        addGenomeToPop(offspring, next, settings) // Adds the genome to the correct species in the new population

    }

    // Release the innovator
    _ = <-inno
    innovate = false
    return

}

func createOffspring(pool GenomeSlice, poolFit float64, everyone GenomeSlice, everyoneFit float64,
    settings *Settings, rng *rand.Rand, inno chan innovator) (offspring *Genome) {
    fmt.Println("createOffspring")
    // Determine the parents
    var p1, p2 *Genome
    p1 = tournament(pool, poolFit, rng)
    if rng.Float64() < settings.InterspeciesMating {
        p2 = tournament(everyone, everyoneFit, rng)
    } else {
        p2 = tournament(pool, poolFit, rng)
    }

    // Order the parents by fitness
    if p1.Fitness[0] < p2.Fitness[0] {
        p1, p2 = p2, p1
    }

    // Create the offspring
    offspring = crossover(p1, p2, settings, rng)

    // Mutate the offspring
    mutateGenome(offspring, rng, settings, inno)

    return
}

func tournament(pool GenomeSlice, poolFit float64, rng *rand.Rand) (gen *Genome) {

    t := rng.Float64() * poolFit
    s := float64(0)
    for i, x := range pool {
        s += x.Fitness[0]
        if s > t {
            return pool[i]
        }
    }

    panic("Should not get here!")
}

// Create the new offspring through crossover. Parent p1 is the more fit
// parent. Matching Genes are inherited randomly, whereas disjoint and
// excess Genes are inherited from the more fit parent. If parents are
// equally fit, then disjoint and excess Genes are also inherited randomly.
// The disabled genes may become enabled again in future generations.
// TODO: Need to check that all NodeGenes used in Connections are present
func crossover(p1, p2 *Genome, settings *Settings, rng *rand.Rand) (child *Genome) {
    fmt.Println("Crossover", p1.ID, p2.ID)

    // Create the new genome
    child = &Genome{ID: settings.ids.Next(),
        NodeGenes: make([]NodeGene, 0, 10),
        ConnGenes: make([]ConnGene, 0, 10)}

    // Determine if both parent's excess and disjoint Genes will be used
    both := (p1.Fitness[0] == p2.Fitness[0])

    // Iterate the NodeGenes
    i, j := 0, 0
    for i < len(p1.NodeGenes) || j < len(p2.NodeGenes) {

        switch {

        case p1.NodeGenes[i].Marker == p2.NodeGenes[j].Marker:
            // Match

            if rng.Float64() > 0.5 {
                child.addNodeGene(p1.NodeGenes[i])
            } else {
                child.addNodeGene(p2.NodeGenes[j])
            }
            i++
            j++

        case j >= len(p2.NodeGenes) || p1.NodeGenes[i].Marker < p2.NodeGenes[j].Marker:
            // Excess or disjoint in p1

            if both {
                if rng.Float64() > 0.5 {
                    child.addNodeGene(p1.NodeGenes[i])
                }
            } else {
                child.addNodeGene(p1.NodeGenes[i])
            }
            i++

        case i >= len(p1.NodeGenes) || p1.NodeGenes[i].Marker > p2.NodeGenes[j].Marker:
            // Excess or disjoint in p2

            if both && rng.Float64() > 0.5 {
                child.addNodeGene(p2.NodeGenes[j])
            }
            j++

        }
    }

    // Iterate the ConnGenes
    i, j = 0, 0
    for i < len(p1.ConnGenes) || j < len(p2.ConnGenes) {

        switch {

        case p1.ConnGenes[i].Marker == p2.ConnGenes[j].Marker:
            // Match

            if rng.Float64() > 0.5 {
                child.addConnGene(p1.ConnGenes[i])
            } else {
                child.addConnGene(p2.ConnGenes[j])
            }
            i++
            j++

        case j >= len(p2.ConnGenes) || p1.ConnGenes[i].Marker < p2.ConnGenes[j].Marker:
            // Excess or disjoint in p1

            if both {
                if rng.Float64() > 0.5 {
                    child.addConnGene(p1.ConnGenes[i])
                }
            } else {
                child.addConnGene(p1.ConnGenes[i])
            }
            i++

        case i >= len(p1.ConnGenes) || p1.ConnGenes[i].Marker > p2.ConnGenes[j].Marker:
            // Excess or disjoint in p2

            if both && rng.Float64() > 0.5 {
                child.addConnGene(p2.ConnGenes[j])
            }
            j++

        }
    }

    return
}

// Returns a rank-ordered slice of Genomes from the organisms passed. Also
// returns the total fitness of the Population. These two things allow us to
// select another Organism via tournament.
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

// Adds a Genome to the Population by finding a compatible Species. If no Species
// is found, a new one will be created with this Genome as its representative. The
// Species in the Population should be sorted by fitness before calling this
// function.
func addGenomeToPop(gen *Genome, pop *Population, settings *Settings) {

    // Look for an existing Species to hold this Genome
    var parent *Species
    for i, x := range pop.Species {

        // Note the compatiblity distance between the Genome and the
        // Species's representative
        d := distance(gen, &x.example, settings)
        if d < settings.CompatThreshold {
            parent = pop.Species[i]
        }

    }

    // No parent found, create a new Species
    if parent == nil {
        parent = &Species{Genomes: []*Genome{}, example: *gen}
        pop.Species = append(pop.Species, parent)
    }
    // Add the Genome to the Species
    parent.Genomes = append(parent.Genomes, gen)

}

// Create the next generation's species
func rollSpecies(curr *Species, settings *Settings, rng *rand.Rand) (next *Species, err error) {

    // Determine the current Species's fitness
    speciesFit(curr, settings)

    // Reverse sort the genomes by fitness
    sort.Sort(curr.Genomes)
    sort.Reverse(curr.Genomes)

    // Remove the weaker members
    survive := int(math.Ceil(float64(len(curr.Genomes)) * settings.SurvivalPercent))
    curr.Genomes = curr.Genomes[0 : survive-1]

    // Create the next species
    next = &Species{Genomes: make([]*Genome, 1, 10), LastUpdate: curr.LastUpdate}

    // Clone the elite genomes
    for i := 0; i < settings.EliteCount && i < len(curr.Genomes); i++ {
        elite := cloneGenome(curr.Genomes[i], curr.Genomes[i].ID)
        next.Genomes = append(next.Genomes, elite)
    }

    // Set the representative
    next.example = *curr.Genomes[rand.Intn(len(curr.Genomes))]

    return

}

// Calculate the fitness of the species by first calculating the adjusted
// fitness of its genomes
func speciesFit(species *Species, settings *Settings) {

    // Special case: only one genome in species
    if len(species.Genomes) == 1 {
        species.fitness = species.Genomes[0].Fitness[0]
        return
    }

    // Determine the adjusted fitness of each genome
    fitness := 0.0
    for i, g := range species.Genomes {

        // Compare this genome to the others in the species
        sum := float64(0)
        for j, x := range species.Genomes {
            if i != j {
                dist := distance(g, x, settings)
                sum += dist
            }
        }

        // Adjust this genome's fitness
        species.Genomes[i].Fitness[0] = sum / float64(len(species.Genomes)-1)

        // Update the species's fitness
        fitness += species.Genomes[i].Fitness[0]

    }

    // Update the species's fitness
    if species.fitness == fitness {
        species.LastUpdate++
    } else {
        species.LastUpdate = 0
        species.fitness = fitness
    }
}

// Calcualtes the compatiblity distance between two genomes. c1, c2 and
// c3 refer to the coefficients described in section 3.3 of "Evolving
// Neural Networks through Augmented Topologies" (Stanley et. al., P109)
func distance(genome1, genome2 *Genome, settings *Settings) float64 {

    // The working variables
    var e, d, n, w, s float64

    // Identify the maximum marker in second genome to determine
    // which genes are disjoint and which are excess
    m := genome2.maxMarker()

    // Process the node genes
    for _, g1 := range genome1.NodeGenes {
        found := false
        for _, g2 := range genome2.NodeGenes {
            if g1.Marker == g2.Marker {
                found = true
                break
            }
        }
        if !found {
            if g1.Marker < m {
                d++ // disjoint
            } else {
                e++ // excess
            }
        }
    }

    // Process the connection genes
    for _, g1 := range genome1.ConnGenes {
        found := false
        for _, g2 := range genome2.ConnGenes {
            if g1.Marker == g2.Marker {
                found = true
                w += math.Abs(g1.Weight - g2.Weight)
                s += 1.0
                break
            }
        }
        if !found {
            if g1.Marker < m {
                d++ // disjoint
            } else {
                e++ // excess
            }
        }
    }

    // Calculate the distance between the two genomes
    c1 := settings.ExcessCoefficient
    c2 := settings.DisjointCoefficient
    c3 := settings.WeightCoefficient
    a := len(genome1.NodeGenes) + len(genome1.ConnGenes)
    b := len(genome2.NodeGenes) + len(genome2.ConnGenes)
    n = math.Max(float64(a), float64(b))
    w = w / s
    return c1*e/n + c2*d/n + c3*w
}
