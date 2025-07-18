// Package hedgehog provides property-based testing for Go.
//
// Hedgehog is a property-based testing library that allows you to write tests
// that are verified against random inputs. It automatically finds minimal
// counterexamples when properties fail.
package hedgehog

import (
	"math/rand"
	"time"
)

// Config holds configuration for property testing.
type Config struct {
	// TestCount is the number of tests to run.
	TestCount int
	// ShrinkCount is the maximum number of shrink attempts.
	ShrinkCount int
	// Seed is the random seed for reproducible tests.
	Seed int64
	// Verbose controls output verbosity.
	Verbose bool
}

// DefaultConfig returns a default configuration for property testing.
func DefaultConfig() *Config {
	return &Config{
		TestCount:   100,
		ShrinkCount: 100,
		Seed:        time.Now().UnixNano(),
		Verbose:     false,
	}
}

// TestResult represents the result of running a property test.
type TestResult int

const (
	// TestPass indicates the property passed all tests.
	TestPass TestResult = iota
	// TestFail indicates the property failed with a counterexample.
	TestFail
	// TestGaveUp indicates the test gave up due to too many discarded cases.
	TestGaveUp
)

// String returns a string representation of the test result.
func (r TestResult) String() string {
	switch r {
	case TestPass:
		return "PASS"
	case TestFail:
		return "FAIL"
	case TestGaveUp:
		return "GAVE UP"
	default:
		return "UNKNOWN"
	}
}

// Random wraps a random number generator for reproducible testing.
type Random struct {
	*rand.Rand
}

// NewRandom creates a new Random instance with the given seed.
func NewRandom(seed int64) *Random {
	return &Random{rand.New(rand.NewSource(seed))}
}

// Split creates a new Random instance from the current one.
func (r *Random) Split() *Random {
	return NewRandom(r.Int63())
}