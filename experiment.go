package neat



type Experiment struct {

    // Settings for the experiement
    Settings Settings

    // The current population - replaced each generation
    Current *Population

    // Helpers
    decoder Decoder
    popEval PopulationEvaluator
	orgEval OrganismEvaluator
	
}

// Sets the Experiment's Decoder. This should be called after loading the
// Experiment from the archive
func (exp *Experiment) SetDecoder(decoder Decoder) {
	exp.decoder = decoder
}

// Sets the Experiment's PopulationEvaluator. This should be called
// after loading the Experiment from the archive. 
func (exp *Experiment) SetPopEval(popEval PopulationEvaluator) {
	exp.popEval = popEval
}

// Closes down the exp by closing open items like sequences
func (exp *Experiment) Close() {

    // Close the sequences
    exp.Settings.Close()

}

