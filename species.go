package neat


// Species is a contatiner of Genomes who share a similar genetic composition.
type Species struct {
    Genomes    GenomeSlice // Collection of Genomes
    LastUpdate uint32      // Number of generations since last fitness change
    fitness    float64     // Adjusted fitness for the Species
    example    Genome      // Representative Genome for comparison
}

// Collection of Species
type SpeciesSlice []*Species

// Len is the number of elements in the collection.
func (s SpeciesSlice) Len() int { return len(s) }

// Less returns whether the element with index i should sort
// before the element with index j.
func (s SpeciesSlice) Less(i, j int) bool {
    // Order Species by their relative fitness
    return s[i].fitness < s[j].fitness
}

// Swap swaps the elements with indexes i and j.
func (s SpeciesSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// Create the next generation's species
func rollSpecies(settings Settings, curr *Species) (next *Species) {

	return
}
