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

type nodeKey struct {
	X, Y float64 // Position of the node in the network
}

type connKey struct {
	Source, Target int // Markers of the source and target nodes
}

type nodeRequest struct {
	key nodeKey
	ret chan int
}

type connRequest struct {
	key connKey
	ret chan int
}

// Innovation provides new IDs and Markers to the different components
// of the NEAT algorithm. The implementation here borrows from the good
// work of others. Specifically, the "sequences" for IDs and Markers are
// based on the AutoInc code found at
//     http://www.mikespook.com/2012/05/golang-funny-play-with-channel/
//
// The maps used to "bless" node and connection genes are based on the
// transactionMap found in David Chsinall's "The Go Programming Language
// Phrasebook" in chapter 10, Concurrency Design Patterns.
type innovation struct {
	running bool // Flag indicating that innovation is still running

	ids     chan int // queue of next available IDs
	markers chan int // queue of next available markers

	nodes map[nodeKey]int // history of node innovations
	conns map[connKey]int // history of connection innovations

	reqN chan nodeRequest
	reqC chan connRequest
}

func newInnovation(pop *Population) *innovation {

	// Create a new innovation
	inno := &innovation{
		running: true,
		ids:     make(chan int, 8),
		markers: make(chan int, 8),

		reqN:  make(chan nodeRequest),
		reqC:  make(chan connRequest),
		nodes: make(map[nodeKey]int),
		conns: make(map[connKey]int)}

	// Identify and start the sequences
	id := 0
	marker := 0
	if pop != nil {
		for _, s := range pop.Species {
			if s.ID > id {
				id = s.ID
			}
			for _, o := range s.Orgs {
				if o.ID > id {
					id = o.ID
				}
				for _, g := range o.Nodes {
					if g.Marker > marker {
						marker = g.Marker
					}
				}
				for _, g := range o.Conns {
					if g.Marker > marker {
						marker = g.Marker
					}
				}
			}
		}
	}
	go inno.startIDs(id + 1)
	go inno.startMarkers(marker + 1)

	// Create and start the marker requests
	inno.reqN = make(chan nodeRequest)
	inno.reqC = make(chan connRequest)
	go inno.runNodes()
	go inno.runConns()

	// Return the innovation
	return inno

}

func (inno *innovation) close() {
	inno.running = false
	close(inno.ids)
	close(inno.markers)
	close(inno.reqC)
	close(inno.reqN)
}

func (inno *innovation) startIDs(start int) {
	defer func() { recover() }()
	for i := start; inno.running; i++ {
		inno.ids <- i
	}
}

func (inno *innovation) startMarkers(start int) {
	defer func() { recover() }()
	for i := start; inno.running; i++ {
		inno.markers <- i
	}
}

func (inno *innovation) nextID() int {
	return <-inno.ids
}

func (inno *innovation) nextMarker() int {
	return <-inno.markers
}

func (inno *innovation) reset() {
	inno.nodes = make(map[nodeKey]int)
	inno.conns = make(map[connKey]int)
}

func (inno *innovation) runNodes() {
	for inno.running {
		req := <-inno.reqN
		m, ok := inno.nodes[req.key]
		if !ok {
			m = <-inno.markers
			inno.nodes[req.key] = m
		}
		req.ret <- m
	}
}

func (inno *innovation) runConns() {
	for inno.running {
		req := <-inno.reqC
		m, ok := inno.conns[req.key]
		if !ok {
			m = <-inno.markers
			inno.conns[req.key] = m
		}
		req.ret <- m
	}
}

func (inno *innovation) blessNodeGene(key nodeKey) int {
	result := make(chan int)
	inno.reqN <- nodeRequest{key, result}
	return <-result
}

func (inno *innovation) blessConnGene(key connKey) int {
	result := make(chan int)
	inno.reqC <- connRequest{key, result}
	return <-result
}
