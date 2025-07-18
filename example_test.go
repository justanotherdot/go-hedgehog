package hedgehog_test

import (
	"slices"
	"strings"
	"testing"
	
	"github.com/rjs/go-hedgehog"
)

// reverse reverses a slice.
func reverse[T any](s []T) []T {
	result := make([]T, len(s))
	for i, v := range s {
		result[len(s)-1-i] = v
	}
	return result
}

// TestReverse demonstrates basic property testing.
func TestReverse(t *testing.T) {
	gen := hedgehog.SliceOf(hedgehog.IntRange(1, 100))
	prop := hedgehog.ForAllNamed(gen, "xs", func(xs []int) bool {
		reversed := reverse(xs)
		doubleReversed := reverse(reversed)
		return slices.Equal(xs, doubleReversed)
	})
	
	report := prop.Run(hedgehog.DefaultConfig())
	report.Assert(t)
}

// TestStringLength demonstrates string property testing.
func TestStringLength(t *testing.T) {
	gen := hedgehog.String()
	prop := hedgehog.ForAllNamed(gen, "text", func(text string) bool {
		uppercase := strings.ToUpper(text)
		return len(uppercase) == len(text)
	})
	
	report := prop.Run(hedgehog.DefaultConfig())
	report.Assert(t)
}

// TestIntegerArithmetic demonstrates integer property testing.
func TestIntegerArithmetic(t *testing.T) {
	gen := hedgehog.IntRange(-1000, 1000)
	prop := hedgehog.ForAllNamed(gen, "x", func(x int) bool {
		return x + 0 == x
	})
	
	report := prop.Run(hedgehog.DefaultConfig())
	report.Assert(t)
}

// TestFailingProperty demonstrates what happens when a property fails.
func TestFailingProperty(t *testing.T) {
	// This property will fail to demonstrate shrinking
	gen := hedgehog.SliceOf(hedgehog.IntRange(1, 100))
	prop := hedgehog.ForAllNamed(gen, "xs", func(xs []int) bool {
		// This is intentionally wrong to show failure
		return len(xs) < 3
	})
	
	config := hedgehog.DefaultConfig()
	config.TestCount = 10 // Run fewer tests for this example
	
	report := prop.Run(config)
	
	if report.Success() {
		t.Log("Property unexpectedly passed")
	} else {
		t.Logf("Property failed as expected:\n%s", report.String())
	}
}

// TestWeightedGeneration demonstrates frequency-based generation.
func TestWeightedGeneration(t *testing.T) {
	gen := hedgehog.Frequency(
		hedgehog.NewWeightedChoice(70, hedgehog.Constant("common")),
		hedgehog.NewWeightedChoice(30, hedgehog.Constant("rare")),
	)
	
	prop := hedgehog.ForAllNamed(gen, "value", func(value string) bool {
		return value == "common" || value == "rare"
	})
	
	report := prop.Run(hedgehog.DefaultConfig())
	report.Assert(t)
}

// TestOneOfGeneration demonstrates uniform choice generation.
func TestOneOfGeneration(t *testing.T) {
	gen := hedgehog.OneOf(
		hedgehog.Constant("red"),
		hedgehog.Constant("green"),
		hedgehog.Constant("blue"),
	)
	
	prop := hedgehog.ForAllNamed(gen, "color", func(color string) bool {
		return color == "red" || color == "green" || color == "blue"
	})
	
	report := prop.Run(hedgehog.DefaultConfig())
	report.Assert(t)
}