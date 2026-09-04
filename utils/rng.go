// Package utils holds small generic utilities that have no domain knowledge.
package utils

import (
	"math/rand"
)

// RNG is a deterministic pseudo-random number generator. All randomness in the
// simulation should go through a World-owned RNG so replays and tests are
// reproducible.
type RNG struct {
	src rand.Source
	rng *rand.Rand
}

// NewRNG returns an RNG seeded with seed.
func NewRNG(seed int64) *RNG {
	src := rand.NewSource(seed)
	return &RNG{src: src, rng: rand.New(src)}
}

// Float64 returns a pseudo-random float64 in [0.0, 1.0).
func (r *RNG) Float64() float64 { return r.rng.Float64() }

// Intn returns a pseudo-random int in [0, n).
func (r *RNG) Intn(n int) int { return r.rng.Intn(n) }

// Norm returns a pseudo-random float64 in [-1.0, 1.0).
func (r *RNG) Norm() float64 { return r.rng.Float64()*2 - 1 }
