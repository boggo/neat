package neat

type Experiment struct {

    // Settings for the experiement
    Settings *Settings

    // The current population - replaced each generation
    Current *Population

    // Helpers
    decoder Decoder
    popEval PopulationEvaluator
    orgEval OrganismEvaluator
}

// Creates a new Experiment from the Settings and data stored in the archive. 
// The Experiment is initialized with a default Decoder and PopulationEvaluator
// which can be overridden with setter methods. The user is responsible for 
// setting the OrganismEvaluator.
func NewExperiment(arch Archiver) (exp *Experiment, err error) {

    // Restore (or load if this is the first time) the Experiment from the
    // Archive
    exp, err = arch.Restore()
    if err != nil {
        return
    }

    // Set the default helpers. These can be overwridden later by calling the
    // appropriate setter methods
    exp.SetDecoder(NewDecoder())
    exp.SetPopEval(NewPopEval())

    // Initialize the sequences in Settings
    initSeq(exp)

	return
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

// Sets the Experiment's PopulationEvaluator. This should be called
// after loading the Experiment from the archive.
func (exp *Experiment) SetOrgEval(orgEval OrganismEvaluator) {
    exp.orgEval = orgEval
}

// Closes down the exp by closing open items like sequences
func (exp *Experiment) Close() {

    // Close the sequences
    exp.Settings.Close()

}
