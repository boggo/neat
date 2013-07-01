package main

import (
    "github.com/boggo/neat"
)

type XOREval struct{}

func (eval XOREval) Evaluate(org *neat.Organism) (fitness []float64, err error) {

    var inputs, outputs []float64
    sum := float64(0)

    // 0 and 0 = 0. Should be less than 0.5. 0 would be perfect answer
    inputs = []float64{0, 0}
    outputs, err = org.Analyze(inputs)
    sum += (0.0 - outputs[0]) * (0.0 - outputs[0])

    // 1 and 1 = 0. Should be less than 0.5. 0 would be perfect answer
    inputs = []float64{1, 1}
    outputs, err = org.Analyze(inputs)
    sum += (0.0 - outputs[0]) * (0.0 - outputs[0])

    // 0 and 1 = 1. Should be less than 0.5. 0 would be perfect answer
    inputs = []float64{1, 0}
    outputs, err = org.Analyze(inputs)
    sum += (1.0 - outputs[0]) * (1.0 - outputs[0])

    // 1 and 0 = 1. Should be less than 0.5. 0 would be perfect answer
    inputs = []float64{0, 1}
    outputs, err = org.Analyze(inputs)
    sum += (1.0 - outputs[0]) * (1.0 - outputs[0])

    fitness = []float64{sum}
	return
	
}

func createFirst() {
    var arch neat.Archiver
    var err error
    arch, err = neat.NewArchiver("/tmp", "xor")
    if err != nil {
        panic(err)
    }
	
	exp := &neat.Experiment{}
	exp.Settings = &neat.Settings{}
	exp.Current = &neat.Population{}
	err = arch.Archive(exp)
	
	
}

func main() {
	//createFirst()
	runOne()
}

func runOne() {

    var arch neat.Archiver
    var err error
    arch, err = neat.NewArchiver("/tmp", "xor")
    if err != nil {
        panic(err)
    }

    var exp *neat.Experiment
    exp, err = neat.NewExperiment(arch)
    if err != nil {
        panic(err)
    }
    exp.SetOrgEval(&XOREval{})

    err = neat.Run(exp, arch, 5)
}
