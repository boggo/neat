package neat

type Settings struct {

    // The number of invididuals in the population
    PopulationSize int

    // Size of the initial genome
    BiasCount   int
    InputCount  int
    OutputCount int

    // Coefficients for calculating distance between genomes
    ExcessCoefficient   float64
    DisjointCoefficient float64
    WeightCoefficient   float64

    // Probabilities for mutation
    MutateWeight        float64
    MutateWeightNew     float64
    MutateEnabled       float64
    MutateAddConnection float64
    MutateAddNode       float64

    // Crossover and breeding probabilities
    Crossover          float64
    InterspeciesMating float64
    AgeToStagnation    int32
    SurvivalPercent    float64 // Percent of a species to survive for mating
    EliteCount         int32   // Number within a species to survive into the next generation
	
	// Prototypical Genome used to track all changes to structure
	ProtoGenome	Genome
	
	// Run-time sequences. These are set during the loading of the Experiment
    ids    *sequence.UInt32
    marker *sequence.UInt32
	
}

// Closes down the exp by closing open items like sequences
func (s *Settings) Close() {

    // Close the sequences
    s.ids.Close()
    s.marker.Close()

}
