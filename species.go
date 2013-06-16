package neat

import (
    "math"
)

// Species is a contatiner of Genomes who share a similar genetic composition.
type Species struct {
    Genomes    GenomeSlice // Collection of Genomes
    LastUpdate int32       // Number of generations since last fitness change
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

// Calcualtes the compatiblity distance between two genomes. c1, c2 and
// c3 refer to the coefficients described in section 3.3 of "Evolving
// Neural Networks through Augmented Topologies" (Stanley et. al., P109)
func distance(genome1, genome2 *Genome, settings *Settings) float64 {

    // The working variables
    var e, d, n, w, s float64

    // Identify the maximum marker in second genome to determine
    // which genes are disjoint and which are excess
    m := maxMarker(genome2)

    // Process the node genes
    for _, g1 := range genome1.NodeGenes {
        found := false
        for _, g2 := range genome2.NodeGenes {
            if g1.Marker == g2.Marker {
                found = true
                break
            }
        }
        if !found {
            if g1.Marker < m {
                d++ // disjoint
            } else {
                e++ // excess
            }
        }
    }

    // Process the connection genes
    for _, g1 := range genome1.ConnGenes {
        found := false
        for _, g2 := range genome2.ConnGenes {
            if g1.Marker == g2.Marker {
                found = true
                w += math.Abs(g1.Weight - g2.Weight)
                s += 1.0
                break
            }
        }
        if !found {
            if g1.Marker < m {
                d++ // disjoint
            } else {
                e++ // excess
            }
        }
    }

    // Calculate the distance between the two genomes
    c1 := settings.ExcessCoefficient
    c2 := settings.DisjointCoefficient
    c3 := settings.WeightCoefficient
    a := len(genome1.NodeGenes) + len(genome1.ConnGenes)
    b := len(genome2.NodeGenes) + len(genome2.ConnGenes)
    n = math.Max(float64(a), float64(b))
    w = w / s
    return c1*e/n + c2*d/n + c3*w
}

// Determine the maximum Marker within a genome
func maxMarker(g *Genome) uint32 {
    m := uint32(0)
    for _, x := range g.NodeGenes {
        if m < x.Marker {
            m = x.Marker
        }
    }
    for _, x := range g.ConnGenes {
        if m < x.Marker {
            m = x.Marker
        }
    }
    return m
}
