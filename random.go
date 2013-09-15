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

package neat

import (
	"math"
	"math/rand"
	"time"
)

type rng struct {
	*rand.Rand
}

var (
	random rng
)

func init() {
	random = rng{rand.New(rand.NewSource(time.Now().UnixNano()))}
}

func (r *rng) Between(a, b float64) float64 {
	return r.Float64()*(b-a) + a
}

func (r *rng) Next() float64 {
	return r.Float64()
}

func (r *rng) Int(n int) int {
	return r.Intn(n)
}

// Returns a normally distributed deviate with zero mean and unit variance.
// From Numerical Recipes in C.
// TODO: involve the mu and sigma parameters. current use mu=0 and sigma=1
var (
	iset bool
	gset float64
)

func (r *rng) Gaussian() float64 {
	var fac, rsq, v1, v2 float64
	if iset == false {
		rsq = 0
		for rsq >= 1.0 || rsq == 0.0 {
			v1 = 2.0*r.Next() - 1.0
			v2 = 2.0*r.Next() - 1.0
			rsq = v1*v1 + v2*v2
		}
		fac = math.Sqrt(-2.0 * math.Log(rsq) / rsq)
		gset = v1 * fac
		iset = true
		return v2 * fac
	} else {
		iset = false
		return gset
	}
}
