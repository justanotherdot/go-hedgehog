package hedgehog

import (
	"testing"
)

func TestConstant(t *testing.T) {
	gen := Constant(42)
	random := NewRandom(12345)

	tree := gen.Generate(random)
	if tree.Value() != 42 {
		t.Errorf("Expected 42, got %d", tree.Value())
	}

	// Constant should have no shrinks
	shrinks := tree.Shrinks()
	if len(shrinks) != 0 {
		t.Errorf("Expected no shrinks, got %d", len(shrinks))
	}
}

func TestIntRange(t *testing.T) {
	gen := IntRange(1, 10)
	random := NewRandom(12345)

	// Generate many values to check range
	for i := 0; i < 100; i++ {
		tree := gen.Generate(random)
		value := tree.Value()

		if value < 1 || value > 10 {
			t.Errorf("Value %d out of range [1, 10]", value)
		}
	}
}

func TestIntRangeShrinking(t *testing.T) {
	gen := IntRange(1, 10)
	random := NewRandom(12345)

	tree := gen.Generate(random)
	value := tree.Value()

	if value == 1 {
		// Skip test if we got the minimum value
		return
	}

	shrinks := tree.Shrinks()
	if len(shrinks) == 0 {
		t.Error("Expected shrinks for non-minimum value")
	}

	// Check that shrinks are smaller
	for _, shrink := range shrinks {
		shrinkValue := shrink.Value()
		if shrinkValue >= value {
			t.Errorf("Shrink %d should be smaller than original %d", shrinkValue, value)
		}
		if shrinkValue < 1 || shrinkValue > 10 {
			t.Errorf("Shrink %d out of range [1, 10]", shrinkValue)
		}
	}
}

func TestBool(t *testing.T) {
	gen := Bool()
	random := NewRandom(12345)

	var trueCount, falseCount int

	// Generate many values to check distribution
	for i := 0; i < 1000; i++ {
		tree := gen.Generate(random)
		if tree.Value() {
			trueCount++
		} else {
			falseCount++
		}
	}

	// Should have some of each (not a strict test, but reasonable)
	if trueCount == 0 || falseCount == 0 {
		t.Errorf("Expected both true and false values, got %d true, %d false", trueCount, falseCount)
	}
}

func TestBoolShrinking(t *testing.T) {
	gen := Bool()
	random := NewRandom(12345)

	// Generate until we get true
	for i := 0; i < 100; i++ {
		tree := gen.Generate(random)
		if tree.Value() {
			// True should shrink to false
			shrinks := tree.Shrinks()
			if len(shrinks) != 1 {
				t.Errorf("Expected 1 shrink for true, got %d", len(shrinks))
			}
			if shrinks[0].Value() != false {
				t.Error("Expected true to shrink to false")
			}
			return
		}
	}

	t.Error("Never generated true value")
}

func TestStringBasic(t *testing.T) {
	gen := String()
	random := NewRandom(12345)

	tree := gen.Generate(random)
	value := tree.Value()

	// Should be a valid string (no specific requirements)
	if len(value) > 20 {
		t.Errorf("String too long: %d chars", len(value))
	}
}

func TestStringShrinking(t *testing.T) {
	gen := String()
	random := NewRandom(12345)

	// Generate until we get a non-empty string
	for i := 0; i < 100; i++ {
		tree := gen.Generate(random)
		value := tree.Value()

		if len(value) > 0 {
			shrinks := tree.Shrinks()
			if len(shrinks) == 0 {
				t.Error("Expected shrinks for non-empty string")
			}

			// First shrink should be empty string
			if shrinks[0].Value() != "" {
				t.Errorf("Expected first shrink to be empty string, got %q", shrinks[0].Value())
			}

			return
		}
	}
}

func TestSliceOf(t *testing.T) {
	gen := SliceOf(IntRange(1, 10))
	random := NewRandom(12345)

	tree := gen.Generate(random)
	value := tree.Value()

	// Check all elements are in range
	for _, elem := range value {
		if elem < 1 || elem > 10 {
			t.Errorf("Element %d out of range [1, 10]", elem)
		}
	}
}

func TestSliceShrinking(t *testing.T) {
	gen := SliceOf(IntRange(1, 10))
	random := NewRandom(12345)

	// Generate until we get a non-empty slice
	for i := 0; i < 100; i++ {
		tree := gen.Generate(random)
		value := tree.Value()

		if len(value) > 0 {
			shrinks := tree.Shrinks()
			if len(shrinks) == 0 {
				t.Error("Expected shrinks for non-empty slice")
			}

			// First shrink should be empty slice
			if len(shrinks[0].Value()) != 0 {
				t.Errorf("Expected first shrink to be empty slice, got length %d", len(shrinks[0].Value()))
			}

			return
		}
	}
}

func TestOneOf(t *testing.T) {
	gen := OneOf(
		Constant("a"),
		Constant("b"),
		Constant("c"),
	)
	random := NewRandom(12345)

	seen := make(map[string]bool)

	// Generate many values to see all options
	for i := 0; i < 100; i++ {
		tree := gen.Generate(random)
		seen[tree.Value()] = true
	}

	// Should have seen all three values
	if !seen["a"] || !seen["b"] || !seen["c"] {
		t.Errorf("Expected to see all values, got %v", seen)
	}
}

func TestFrequency(t *testing.T) {
	gen := Frequency(
		NewWeightedChoice(90, Constant("common")),
		NewWeightedChoice(10, Constant("rare")),
	)
	random := NewRandom(12345)

	counts := make(map[string]int)

	// Generate many values to check distribution
	for i := 0; i < 1000; i++ {
		tree := gen.Generate(random)
		counts[tree.Value()]++
	}

	// Should see both values
	if counts["common"] == 0 || counts["rare"] == 0 {
		t.Errorf("Expected both values, got %v", counts)
	}

	// Common should be more frequent (not a strict test)
	if counts["common"] < counts["rare"] {
		t.Errorf("Expected 'common' to be more frequent than 'rare', got %v", counts)
	}
}

func TestGenMap(t *testing.T) {
	gen := IntRange(1, 10)
	mapped := gen.Map(func(x int) int { return x * 2 })

	random := NewRandom(12345)
	tree := mapped.Generate(random)
	value := tree.Value()

	// Should be in range [2, 20]
	if value < 2 || value > 20 {
		t.Errorf("Mapped value %d out of expected range [2, 20]", value)
	}

	// Should be even
	if value%2 != 0 {
		t.Errorf("Mapped value %d should be even", value)
	}
}

// Edge cases and error conditions

func TestIntRangeInvalidRange(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for invalid range")
		}
	}()

	IntRange(10, 5) // min > max should panic
}

func TestFloat64RangeInvalidInputs(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for invalid range")
		}
	}()

	Float64Range(10.0, 5.0) // min > max should panic
}

func TestFrequencyValidation(t *testing.T) {
	t.Run("empty choices", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for empty choices")
			}
		}()

		Frequency[string]() // No choices should panic
	})

	t.Run("zero total weight", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for zero total weight")
			}
		}()

		Frequency(
			NewWeightedChoice(0, Constant("zero")),
			NewWeightedChoice(0, Constant("also zero")),
		) // All zero weights should panic
	})

	t.Run("negative weight", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for negative weight")
			}
		}()

		Frequency(
			NewWeightedChoice(-1, Constant("negative")),
			NewWeightedChoice(5, Constant("positive")),
		) // Negative weight should panic
	})
}

func TestOneOfEmptyGenerators(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for empty generators")
		}
	}()

	OneOf[string]() // No generators should panic
}

func TestStringWithZeroLength(t *testing.T) {
	gen := StringWithLength(Constant(0))
	random := NewRandom(12345)

	tree := gen.Generate(random)
	value := tree.Value()

	if value != "" {
		t.Errorf("Expected empty string, got %q", value)
	}

	// Should have no shrinks
	shrinks := tree.Shrinks()
	if len(shrinks) != 0 {
		t.Errorf("Expected no shrinks for empty string, got %d", len(shrinks))
	}
}

func TestSliceWithZeroLength(t *testing.T) {
	gen := SliceOfWithLength(IntRange(1, 10), Constant(0))
	random := NewRandom(12345)

	tree := gen.Generate(random)
	value := tree.Value()

	if len(value) != 0 {
		t.Errorf("Expected empty slice, got length %d", len(value))
	}

	// Should have no shrinks
	shrinks := tree.Shrinks()
	if len(shrinks) != 0 {
		t.Errorf("Expected no shrinks for empty slice, got %d", len(shrinks))
	}
}
