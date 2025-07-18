package hedgehog

import (
	"fmt"
	"math"
)

// Gen represents a generator that produces values of type T along with a tree
// of shrinks for that value.
type Gen[T any] struct {
	generate func(*Random) *Tree[T]
}

// NewGen creates a new generator from a generation function.
func NewGen[T any](generate func(*Random) *Tree[T]) *Gen[T] {
	return &Gen[T]{generate: generate}
}

// Generate produces a value using the given random source.
func (g *Gen[T]) Generate(r *Random) *Tree[T] {
	return g.generate(r)
}

// Map transforms a generator by applying a function to its values.
func (g *Gen[T]) Map(f func(T) T) *Gen[T] {
	return NewGen(func(r *Random) *Tree[T] {
		tree := g.Generate(r)
		return tree.Map(f)
	})
}

// Constant creates a generator that always produces the same value.
func Constant[T any](value T) *Gen[T] {
	return NewGen(func(r *Random) *Tree[T] {
		return NewTree(value, func() []*Tree[T] {
			return nil
		})
	})
}

// IntRange creates a generator for integers in the given range [min, max].
func IntRange(min, max int) *Gen[int] {
	if min > max {
		panic(fmt.Sprintf("IntRange: min (%d) > max (%d)", min, max))
	}

	return NewGen(func(r *Random) *Tree[int] {
		// Calculate range size safely to avoid overflow
		rangeSize := int64(max) - int64(min) + 1
		if rangeSize <= 0 || rangeSize > int64(^uint(0)>>1) {
			// Range too large for safe generation, clamp to reasonable size
			rangeSize = 1000
		}

		value := r.Intn(int(rangeSize)) + min
		return NewTree(value, func() []*Tree[int] {
			return shrinkInt(value, min, max)
		})
	})
}

// shrinkInt generates shrinks for an integer value within bounds.
func shrinkInt(value, min, max int) []*Tree[int] {
	if value == min {
		return nil
	}

	var shrinks []*Tree[int]

	// Try shrinking towards zero (if in bounds)
	if min <= 0 && 0 <= max && value != 0 {
		shrinks = append(shrinks, NewTree(0, func() []*Tree[int] {
			return shrinkInt(0, min, max)
		}))
	}

	// Try shrinking towards min
	if value > min {
		candidate := min + (value-min)/2
		if candidate != value {
			shrinks = append(shrinks, NewTree(candidate, func() []*Tree[int] {
				return shrinkInt(candidate, min, max)
			}))
		}
	}

	return shrinks
}

// Float64Range creates a generator for float64 values in the given range [min, max].
func Float64Range(min, max float64) *Gen[float64] {
	if min > max {
		panic(fmt.Sprintf("Float64Range: min (%f) > max (%f)", min, max))
	}
	if math.IsNaN(min) || math.IsNaN(max) {
		panic("Float64Range: min or max is NaN")
	}
	if math.IsInf(min, 0) || math.IsInf(max, 0) {
		panic("Float64Range: min or max is infinite")
	}

	return NewGen(func(r *Random) *Tree[float64] {
		value := r.Float64()*(max-min) + min
		return NewTree(value, func() []*Tree[float64] {
			return shrinkFloat64(value, min, max)
		})
	})
}

// shrinkFloat64 generates shrinks for a float64 value within bounds.
func shrinkFloat64(value, min, max float64) []*Tree[float64] {
	if math.Abs(value-min) < 1e-10 {
		return nil
	}

	var shrinks []*Tree[float64]

	// Try shrinking towards zero (if in bounds)
	if min <= 0 && 0 <= max && math.Abs(value) > 1e-10 {
		shrinks = append(shrinks, NewTree(0.0, func() []*Tree[float64] {
			return shrinkFloat64(0.0, min, max)
		}))
	}

	// Try shrinking towards min
	if value > min {
		candidate := min + (value-min)/2
		if math.Abs(candidate-value) > 1e-10 {
			shrinks = append(shrinks, NewTree(candidate, func() []*Tree[float64] {
				return shrinkFloat64(candidate, min, max)
			}))
		}
	}

	return shrinks
}

// Bool creates a generator for boolean values.
func Bool() *Gen[bool] {
	return NewGen(func(r *Random) *Tree[bool] {
		value := r.Intn(2) == 1
		return NewTree(value, func() []*Tree[bool] {
			if value {
				return []*Tree[bool]{
					NewTree(false, func() []*Tree[bool] { return nil }),
				}
			}
			return nil
		})
	})
}

// String creates a generator for strings of ASCII characters.
func String() *Gen[string] {
	return StringWithLength(IntRange(0, 20))
}

// StringWithLength creates a generator for strings with the given length generator.
func StringWithLength(lengthGen *Gen[int]) *Gen[string] {
	return NewGen(func(r *Random) *Tree[string] {
		lengthTree := lengthGen.Generate(r)
		length := lengthTree.Value()

		if length == 0 {
			return NewTree("", func() []*Tree[string] {
				return nil
			})
		}

		var chars []rune
		for i := 0; i < length; i++ {
			// Generate printable ASCII characters
			char := rune(r.Intn(95) + 32)
			chars = append(chars, char)
		}

		value := string(chars)
		return NewTree(value, func() []*Tree[string] {
			return shrinkString(value)
		})
	})
}

// shrinkString generates shrinks for a string value.
func shrinkString(value string) []*Tree[string] {
	if len(value) == 0 {
		return nil
	}

	var shrinks []*Tree[string]

	// Try empty string
	shrinks = append(shrinks, NewTree("", func() []*Tree[string] {
		return nil
	}))

	// Try removing characters
	runes := []rune(value)
	for i := 0; i < len(runes); i++ {
		candidate := string(runes[:i]) + string(runes[i+1:])
		if candidate != value {
			shrinks = append(shrinks, NewTree(candidate, func() []*Tree[string] {
				return shrinkString(candidate)
			}))
		}
	}

	// Try simplifying characters
	for i, ch := range runes {
		simplified := simplifyChar(ch)
		if simplified != ch {
			newRunes := make([]rune, len(runes))
			copy(newRunes, runes)
			newRunes[i] = simplified
			candidate := string(newRunes)
			shrinks = append(shrinks, NewTree(candidate, func() []*Tree[string] {
				return shrinkString(candidate)
			}))
		}
	}

	return shrinks
}

// simplifyChar attempts to simplify a character for shrinking.
func simplifyChar(ch rune) rune {
	switch {
	case ch >= 'A' && ch <= 'Z':
		return 'a' + (ch - 'A')
	case ch >= '1' && ch <= '9':
		return '0'
	case ch == '0':
		return '0'
	case ch >= 'b' && ch <= 'z':
		return 'a'
	case ch == 'a':
		return 'a'
	default:
		return 'a'
	}
}

// SliceOf creates a generator for slices of elements generated by the given generator.
func SliceOf[T any](elemGen *Gen[T]) *Gen[[]T] {
	return SliceOfWithLength(elemGen, IntRange(0, 10))
}

// SliceOfWithLength creates a generator for slices with the given length and element generators.
func SliceOfWithLength[T any](elemGen *Gen[T], lengthGen *Gen[int]) *Gen[[]T] {
	return NewGen(func(r *Random) *Tree[[]T] {
		lengthTree := lengthGen.Generate(r)
		length := lengthTree.Value()

		if length == 0 {
			return NewTree([]T{}, func() []*Tree[[]T] {
				return nil
			})
		}

		var elements []T
		for i := 0; i < length; i++ {
			elemTree := elemGen.Generate(r)
			elements = append(elements, elemTree.Value())
		}

		return NewTree(elements, func() []*Tree[[]T] {
			return shrinkSlice(elements, elemGen, r)
		})
	})
}

// shrinkSlice generates shrinks for a slice value.
func shrinkSlice[T any](value []T, elemGen *Gen[T], r *Random) []*Tree[[]T] {
	if len(value) == 0 {
		return nil
	}

	var shrinks []*Tree[[]T]

	// Try empty slice
	shrinks = append(shrinks, NewTree([]T{}, func() []*Tree[[]T] {
		return nil
	}))

	// Try removing elements
	for i := 0; i < len(value); i++ {
		candidate := append(value[:i:i], value[i+1:]...)
		shrinks = append(shrinks, NewTree(candidate, func() []*Tree[[]T] {
			return shrinkSlice(candidate, elemGen, r)
		}))
	}

	return shrinks
}

// OneOf creates a generator that chooses uniformly from the given generators.
func OneOf[T any](generators ...*Gen[T]) *Gen[T] {
	if len(generators) == 0 {
		panic("OneOf: no generators provided")
	}

	return NewGen(func(r *Random) *Tree[T] {
		index := r.Intn(len(generators))
		return generators[index].Generate(r)
	})
}

// Frequency creates a generator that chooses from weighted generators.
func Frequency[T any](choices ...WeightedChoice[T]) *Gen[T] {
	if len(choices) == 0 {
		panic("Frequency: no choices provided")
	}

	totalWeight := 0
	for _, choice := range choices {
		if choice.Weight < 0 {
			panic(fmt.Sprintf("Frequency: negative weight %d", choice.Weight))
		}
		totalWeight += choice.Weight
	}

	if totalWeight == 0 {
		panic("Frequency: total weight is zero")
	}

	return NewGen(func(r *Random) *Tree[T] {
		pick := r.Intn(totalWeight)
		current := 0

		for _, choice := range choices {
			current += choice.Weight
			if pick < current {
				return choice.Generator.Generate(r)
			}
		}

		// Fallback to first choice
		return choices[0].Generator.Generate(r)
	})
}

// WeightedChoice represents a generator with an associated weight.
type WeightedChoice[T any] struct {
	Weight    int
	Generator *Gen[T]
}

// NewWeightedChoice creates a new weighted choice.
func NewWeightedChoice[T any](weight int, generator *Gen[T]) WeightedChoice[T] {
	return WeightedChoice[T]{Weight: weight, Generator: generator}
}
