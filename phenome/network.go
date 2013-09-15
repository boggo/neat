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

package phenome

import (
	"errors"
	"github.com/boggo/neat"
	"github.com/boggo/neural"
)

// Default implementation of Phenome uses a neural.Network to analyze the
// inputs
type networkPhenome struct {
	network *neural.Network
}

func NewNetwork(network *neural.Network) neat.Phenome {
	return &networkPhenome{network}
}

// Analyzes the inputs and returns the results as a slice of float64. This
// implementation is for the networkPhenome. Network does not return an error
// in Activate() so the error value is nil upon return.
func (p *networkPhenome) Analyze(inputs []float64) (outputs []float64, err error) {
	outputs = p.network.Activate(inputs)
	switch {
	case outputs == nil:
		err = errors.New("Network produced nil for outputs")
	case len(outputs) == 0:
		err = errors.New("Network produced outputs with zero length")
	}
	return
}
