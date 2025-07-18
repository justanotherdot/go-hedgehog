package hedgehog

import (
	"fmt"
	"strings"
)

// Tree represents a value with a lazy tree of shrinks.
type Tree[T any] struct {
	value   T
	shrinks func() []*Tree[T]
}

// NewTree creates a new tree with the given value and shrink function.
func NewTree[T any](value T, shrinks func() []*Tree[T]) *Tree[T] {
	return &Tree[T]{
		value:   value,
		shrinks: shrinks,
	}
}

// Value returns the value at the root of the tree.
func (t *Tree[T]) Value() T {
	return t.value
}

// Shrinks returns the immediate shrinks of this tree.
func (t *Tree[T]) Shrinks() []*Tree[T] {
	if t.shrinks == nil {
		return nil
	}
	return t.shrinks()
}

// Map applies a function to the value and all shrinks in the tree.
func (t *Tree[T]) Map(f func(T) T) *Tree[T] {
	return NewTree(f(t.value), func() []*Tree[T] {
		shrinks := t.Shrinks()
		if len(shrinks) == 0 {
			return nil
		}
		
		var mappedShrinks []*Tree[T]
		for _, shrink := range shrinks {
			mappedShrinks = append(mappedShrinks, shrink.Map(f))
		}
		return mappedShrinks
	})
}

// Unfold creates a tree by repeatedly applying shrink functions.
func (t *Tree[T]) Unfold() []*Tree[T] {
	var result []*Tree[T]
	var queue []*Tree[T]
	
	queue = append(queue, t)
	
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		
		result = append(result, current)
		
		shrinks := current.Shrinks()
		queue = append(queue, shrinks...)
	}
	
	return result
}

// Render creates a string representation of the tree for debugging.
func (t *Tree[T]) Render() string {
	return t.renderWithDepth(0, 3)
}

// renderWithDepth renders the tree with indentation up to a certain depth.
func (t *Tree[T]) renderWithDepth(depth, maxDepth int) string {
	if depth > maxDepth {
		return ""
	}
	
	indent := strings.Repeat("  ", depth)
	result := fmt.Sprintf("%s%v\n", indent, t.value)
	
	shrinks := t.Shrinks()
	for _, shrink := range shrinks {
		result += shrink.renderWithDepth(depth+1, maxDepth)
	}
	
	return result
}

// RenderCompact creates a compact string representation showing just the values.
func (t *Tree[T]) RenderCompact() string {
	var values []string
	trees := t.Unfold()
	
	for _, tree := range trees {
		values = append(values, fmt.Sprintf("%v", tree.Value()))
	}
	
	return strings.Join(values, " -> ")
}

// Flatten returns all values in the tree in breadth-first order.
func (t *Tree[T]) Flatten() []T {
	var result []T
	trees := t.Unfold()
	
	for _, tree := range trees {
		result = append(result, tree.Value())
	}
	
	return result
}

// Filter returns a new tree containing only shrinks that satisfy the predicate.
func (t *Tree[T]) Filter(pred func(T) bool) *Tree[T] {
	return NewTree(t.value, func() []*Tree[T] {
		shrinks := t.Shrinks()
		if len(shrinks) == 0 {
			return nil
		}
		
		var filteredShrinks []*Tree[T]
		for _, shrink := range shrinks {
			if pred(shrink.Value()) {
				filteredShrinks = append(filteredShrinks, shrink.Filter(pred))
			}
		}
		return filteredShrinks
	})
}

// Prune limits the tree to a maximum depth.
func (t *Tree[T]) Prune(maxDepth int) *Tree[T] {
	if maxDepth <= 0 {
		return NewTree(t.value, func() []*Tree[T] { return nil })
	}
	
	return NewTree(t.value, func() []*Tree[T] {
		shrinks := t.Shrinks()
		if len(shrinks) == 0 {
			return nil
		}
		
		var prunedShrinks []*Tree[T]
		for _, shrink := range shrinks {
			prunedShrinks = append(prunedShrinks, shrink.Prune(maxDepth-1))
		}
		return prunedShrinks
	})
}

// FindSmallest finds the smallest value in the tree that satisfies the predicate.
func (t *Tree[T]) FindSmallest(pred func(T) bool) (T, bool) {
	if !pred(t.value) {
		var zero T
		return zero, false
	}
	
	// Try to find a smaller value among the shrinks
	for _, shrink := range t.Shrinks() {
		if smallest, found := shrink.FindSmallest(pred); found {
			return smallest, true
		}
	}
	
	// If no smaller value found, return the current value
	return t.value, true
}