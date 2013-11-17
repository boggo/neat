/*  Copyright (c) 2013, Brian Hummer (brian@boggo.net)
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:
    * Redistributions of source code must retain the above copyright
      notice, this list of conditions and the following disclaimer.
    * Redistributions in binary form must reproduce the above copyright
      notice, this list of conditions and the following disclaimer in the
      documentation and/or other materials provided with the distribution.
    * Neither the name of the boggo.net nor the
      names of its contributors may be used to endorse or promote products
      derived from this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL BRIAN HUMMER BE LIABLE FOR ANY
DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
(INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

package reporter

import (
	"fmt"
	"github.com/boggo/neat"
	"math"
	"strconv"
	"strings"
	"time"
)

type consoleReporter struct{}

func NewConsole() *consoleReporter {
	return &consoleReporter{}
}

func (rep *consoleReporter) Report(pop *neat.Population) (err error) {

	// Report the progress
	fmt.Println("=============================================================================")
	fmt.Println("Generation", pop.Generation, "finished", time.Now())

	fmt.Println("-----------------------------------------------------------------------------")
	n := len(pop.Species)
	fmt.Println(n, "Species")
	fmt.Printf("  ID    Count   Fitness   Nodes   Conns    Age    Stagn\n")
	fmt.Printf("------ ------- --------- ------- ------- ------- -------\n")

	var b float64
	var m, l int
	l = math.MaxInt32
	var bs, ms, ls *neat.Organism

	fit := make([]float64, n)
	nod := make([]float64, n)
	con := make([]float64, n)
	age := make([]float64, n)
	sta := make([]float64, n)
	for i, s := range pop.Species {
		fit[i] = avgFit(s)
		nod[i] = avgNod(s)
		con[i] = avgCon(s)
		age[i] = float64(s.Age)
		sta[i] = float64(s.Age - s.BestFitAge)
		fmt.Printf("%6d %7d %7v %7v %7v %7v %7v\n", s.ID, len(s.Orgs),
			fmtf(fit[i], 5, 9),
			fmtf(nod[i], 2, 7),
			fmtf(con[i], 2, 7),
			fmtf(age[i], 0, 7),
			fmtf(sta[i], 0, 7))

		// Organisms summary
		for _, o := range s.Orgs {

			// Update best
			if o.Fitness[0] > b {
				bs = o
				b = o.Fitness[0]
			}
			c := len(o.Nodes) + len(o.Conns)
			switch {
			case c < l: // least complexity
				ls = o
				l = c
			case c > m: // Max complextity
				ms = o
				m = c
			case c == m && ms != nil: // Max complexity with higher fitness
				if ms.Fitness[0] < o.Fitness[0] {
					ms = o
				}
			case c == l && ls != nil: // Least complexity with higher fitness
				if ls.Fitness[0] < o.Fitness[0] {
					ls = o
				}
			}
		}
	}

	fmt.Println("")
	fmt.Println("Best Fitness: ", bs)
	fmt.Println("Most Complex: ", ms)
	fmt.Println("Least Complex:", ls)
	fmt.Println("-----------------------------------------------------------------------------")

	// Return the error if any
	return

}

func fmtf(f float64, d int, l int) string {
	a := fmt.Sprintf(("%." + strconv.Itoa(d) + "f"), f)
	i := l - len(a)
	if i > 0 {
		a = strings.Repeat(" ", i) + a
	}
	return a
}
