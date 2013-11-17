/*
CREDIT WHERE CREDIT IS DUE:
This experiment is based on the Pole Balancing experiement in neat-python
(https://code.google.com/p/neat-python/), specifically the Evaluate
function below is based on the evaluate_population and cart_pole functions in single_pole.py. That code is GPL3 and copyright belongs to that software's authors.

The remaining code in this experiment was written by Brian Hummer (brian@boggo.net)
and is released under GPL3 because the example used was GPL3. The libraries used,
though, are 3-clause, "new" BSD licensed.

GPL3 license:

Copyright (C) 2013 Brian Hummer (brian@boggo.net)

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package main

import (
	"errors"
	//"fmt"
	"github.com/boggo/neat"
	"github.com/boggo/neat/archiver"
	"github.com/boggo/neat/decoder"
	"github.com/boggo/neat/popeval"
	"github.com/boggo/neat/reporter"
	"github.com/boggo/neat/settings"
	"math"
	//	"math/rand"
	//"time"
)

/* Original eval function from neat-python:

def cart_pole(net_output, x, x_dot, theta, theta_dot):
    ''' Directly copied from Stanley's C++ source code '''

    GRAVITY = 9.8
    MASSCART = 1.0
    MASSPOLE = 0.1
    TOTAL_MASS = (MASSPOLE + MASSCART)
    LENGTH = 0.5    # actually half the pole's length
    POLEMASS_LENGTH = (MASSPOLE * LENGTH)
    FORCE_MAG = 10.0
    TAU = 0.02  # seconds between state updates
    FOURTHIRDS = 1.3333333333333

    #force = (net_output - 0.5) * FORCE_MAG * 2
    if net_output > 0.5:
        force = FORCE_MAG
    else:
        force = -FORCE_MAG

    costheta = math.cos(theta)
    sintheta = math.sin(theta)

    temp = (force + POLEMASS_LENGTH * theta_dot * theta_dot * sintheta)/ TOTAL_MASS

    thetaacc = (GRAVITY*sintheta - costheta*temp)\
               /(LENGTH * (FOURTHIRDS - MASSPOLE * costheta * costheta/TOTAL_MASS))

    xacc  = temp - POLEMASS_LENGTH * thetaacc * costheta / TOTAL_MASS

    #Update the four state variables, using Euler's method
    x         += TAU * x_dot
    x_dot     += TAU * xacc
    theta     += TAU * theta_dot
    theta_dot += TAU * thetaacc

    return x, x_dot, theta, theta_dot

def evaluate_population(population):

    twelve_degrees = 0.2094384 #radians
    num_steps = 10**5

    for chromo in population:

        net = nn.create_phenotype(chromo)

        # initial conditions (as used by Stanley)
        x         = random.randint(0, 4799)/1000.0 - 2.4
        x_dot     = random.randint(0, 1999)/1000.0 - 1.0
        theta     = random.randint(0,  399)/1000.0 - 0.2
        theta_dot = random.randint(0, 2999)/1000.0 - 1.5
        #x = 0.0
        #x_dot = 0.0
        #theta = 0.0
        #theta_dot = 0.0

        fitness = 0

        for trials in xrange(num_steps):

            # maps into [0,1]
            inputs = [(x + 2.4)/4.8,
                      (x_dot + 0.75)/1.5,
                      (theta + twelve_degrees)/0.41,
                      (theta_dot + 1.0)/2.0]

            # a normalizacao so acontece para estas condicoes iniciais
            # nada garante que a evolucao do sistema leve a outros
            # valores de x, x_dot e etc...

            action = net.pactivate(inputs)

            # Apply action to the simulated cart-pole
            x, x_dot, theta, theta_dot = cart_pole(action[0], x, x_dot, theta, theta_dot)

            # Check for failure.  If so, return steps
            # the number of steps indicates the fitness: higher = better
            fitness += 1
            if (abs(x) >= 2.4 or abs(theta) >= twelve_degrees):
            #if abs(theta) > twelve_degrees: # Igel (p. 5) uses theta criteria only
                # the cart/pole has run/inclined out of the limits
                break

        chromo.fitness = fitness

*/

func cartPole(netOutput float64, x, x_dot, theta, theta_dot float64) (float64, float64, float64, float64) {
	//''' Directly copied from Stanley's C++ source code '''

	GRAVITY := float64(9.8)
	MASSCART := float64(1.0)
	MASSPOLE := float64(0.1)
	TOTAL_MASS := (MASSPOLE + MASSCART)
	LENGTH := float64(0.5) // actually half the pole's length
	POLEMASS_LENGTH := (MASSPOLE * LENGTH)
	FORCE_MAG := float64(10.0)
	TAU := float64(0.02) // seconds between state updates
	FOURTHIRDS := float64(1.3333333333333)

	var force float64
	if netOutput > 0.5 {
		force = FORCE_MAG
	} else {
		force = -FORCE_MAG
	}

	costheta := math.Cos(theta)
	sintheta := math.Sin(theta)

	temp := (force + POLEMASS_LENGTH*theta_dot*theta_dot*sintheta) / TOTAL_MASS

	thetaacc := (GRAVITY*sintheta - costheta*temp) / (LENGTH * (FOURTHIRDS - MASSPOLE*costheta*costheta/TOTAL_MASS))

	xacc := temp - POLEMASS_LENGTH*thetaacc*costheta/TOTAL_MASS

	//Update the four state variables, using Euler's method
	x += TAU * x_dot
	x_dot += TAU * xacc
	theta += TAU * theta_dot
	theta_dot += TAU * thetaacc

	return x, x_dot, theta, theta_dot
}

type singPoleEval struct{}

func (eval singPoleEval) Evaluate(org *neat.Organism) (err error) {

	if org.Phenome == nil {
		err = errors.New("Cannot evaluate an org without a Phenome")
		org.Fitness = []float64{0} // Minimal fitness
		return
	}

	twelve_degrees := float64(0.2094384) // radians
	num_steps := int(math.Pow(10, 5))

	// initial conditions (as used by Stanley)
	//rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	//x := float64(rng.Intn(4800))/1000.0 - 2.4
	//x_dot := float64(rng.Intn(2000))/1000.0 - 1.0
	//theta := float64(rng.Intn(400))/1000.0 - 0.2
	//theta_dot := float64(rng.Intn(3000))/1000.0 - 1.5
	var x, x_dot, theta, theta_dot float64

	fitness := float64(0)

	for trials := 0; trials < num_steps; trials++ {

		// maps into [0,1]
		inputs := []float64{(x + 2.4) / 4.8,
			(x_dot + 0.75) / 1.5,
			(theta + twelve_degrees) / 0.41,
			(theta_dot + 1.0) / 2.0}

		// a normalizacao so acontece para estas condicoes iniciais
		// nada garante que a evolucao do sistema leve a outros
		// valores de x, x_dot e etc...

		action, err2 := org.Analyze(inputs)
		if err2 != nil {
			err = err2
			org.Fitness = []float64{0}
			return
		}
		// Apply action to the simulated cart-pole
		x, x_dot, theta, theta_dot = cartPole(action[0], x, x_dot, theta, theta_dot)

		// Check for failure.  If so, return steps
		// the number of steps indicates the fitness: higher = better
		fitness += 1
		//if fitness > 100 {
		//	fmt.Println("WTF?")
		//}
		if math.Abs(x) >= 2.4 || (math.Abs(theta) >= twelve_degrees) {
			break
		}
	}
	org.Fitness = []float64{fitness}
	return
}

func main() {

	// Load the settings
	ldr := settings.NewJSON("singpole-settings.json")
	s, err := ldr.Load()
	if err != nil {
		panic(err)
	}

	// Create the archiver
	a := archiver.NewJSON("singpole-pop.json")

	// Create the reporter
	r := reporter.NewConsole()

	// Create the evaluators
	o := &singPoleEval{}
	p := popeval.NewConcurrent()

	// Create the decoder
	d := decoder.NewNEAT()

	// Iterate the experiment
	neat.Iterate(s, 25, d, p, o, a, r)
}
