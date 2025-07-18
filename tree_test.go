package hedgehog

import (
	"testing"
)

func TestTreeValue(t *testing.T) {
	tree := NewTree(42, func() []*Tree[int] { return nil })

	if tree.Value() != 42 {
		t.Errorf("Expected 42, got %d", tree.Value())
	}
}

func TestTreeShrinks(t *testing.T) {
	tree := NewTree(10, func() []*Tree[int] {
		return []*Tree[int]{
			NewTree(5, func() []*Tree[int] { return nil }),
			NewTree(0, func() []*Tree[int] { return nil }),
		}
	})

	shrinks := tree.Shrinks()
	if len(shrinks) != 2 {
		t.Errorf("Expected 2 shrinks, got %d", len(shrinks))
	}

	if shrinks[0].Value() != 5 {
		t.Errorf("Expected first shrink to be 5, got %d", shrinks[0].Value())
	}

	if shrinks[1].Value() != 0 {
		t.Errorf("Expected second shrink to be 0, got %d", shrinks[1].Value())
	}
}

func TestTreeMap(t *testing.T) {
	tree := NewTree(5, func() []*Tree[int] {
		return []*Tree[int]{
			NewTree(3, func() []*Tree[int] { return nil }),
			NewTree(1, func() []*Tree[int] { return nil }),
		}
	})

	mapped := tree.Map(func(x int) int { return x * 2 })

	if mapped.Value() != 10 {
		t.Errorf("Expected mapped value 10, got %d", mapped.Value())
	}

	shrinks := mapped.Shrinks()
	if len(shrinks) != 2 {
		t.Errorf("Expected 2 mapped shrinks, got %d", len(shrinks))
	}

	if shrinks[0].Value() != 6 {
		t.Errorf("Expected first mapped shrink to be 6, got %d", shrinks[0].Value())
	}

	if shrinks[1].Value() != 2 {
		t.Errorf("Expected second mapped shrink to be 2, got %d", shrinks[1].Value())
	}
}

func TestTreeFlatten(t *testing.T) {
	tree := NewTree(10, func() []*Tree[int] {
		return []*Tree[int]{
			NewTree(5, func() []*Tree[int] {
				return []*Tree[int]{
					NewTree(2, func() []*Tree[int] { return nil }),
				}
			}),
			NewTree(0, func() []*Tree[int] { return nil }),
		}
	})

	values := tree.Flatten()
	expected := []int{10, 5, 0, 2}

	if len(values) != len(expected) {
		t.Errorf("Expected %d values, got %d", len(expected), len(values))
	}

	for i, expected := range expected {
		if i >= len(values) || values[i] != expected {
			t.Errorf("At index %d, expected %d, got %d", i, expected, values[i])
		}
	}
}

func TestTreeRenderCompact(t *testing.T) {
	tree := NewTree(10, func() []*Tree[int] {
		return []*Tree[int]{
			NewTree(5, func() []*Tree[int] { return nil }),
			NewTree(0, func() []*Tree[int] { return nil }),
		}
	})

	rendered := tree.RenderCompact()
	expected := "10 -> 5 -> 0"

	if rendered != expected {
		t.Errorf("Expected %q, got %q", expected, rendered)
	}
}

func TestTreeFilter(t *testing.T) {
	tree := NewTree(10, func() []*Tree[int] {
		return []*Tree[int]{
			NewTree(5, func() []*Tree[int] { return nil }),
			NewTree(3, func() []*Tree[int] { return nil }),
			NewTree(0, func() []*Tree[int] { return nil }),
		}
	})

	// Filter to keep only values > 2
	filtered := tree.Filter(func(x int) bool { return x > 2 })

	if filtered.Value() != 10 {
		t.Errorf("Expected filtered value 10, got %d", filtered.Value())
	}

	shrinks := filtered.Shrinks()
	if len(shrinks) != 2 {
		t.Errorf("Expected 2 filtered shrinks, got %d", len(shrinks))
	}

	// Should have 5 and 3, but not 0
	values := []int{shrinks[0].Value(), shrinks[1].Value()}
	if (values[0] != 5 || values[1] != 3) && (values[0] != 3 || values[1] != 5) {
		t.Errorf("Expected shrinks to be 5 and 3, got %v", values)
	}
}

func TestTreePrune(t *testing.T) {
	tree := NewTree(10, func() []*Tree[int] {
		return []*Tree[int]{
			NewTree(5, func() []*Tree[int] {
				return []*Tree[int]{
					NewTree(2, func() []*Tree[int] {
						return []*Tree[int]{
							NewTree(1, func() []*Tree[int] { return nil }),
						}
					}),
				}
			}),
		}
	})

	// Prune to depth 2
	pruned := tree.Prune(2)

	if pruned.Value() != 10 {
		t.Errorf("Expected pruned value 10, got %d", pruned.Value())
	}

	// Should have first level shrinks
	shrinks := pruned.Shrinks()
	if len(shrinks) != 1 {
		t.Errorf("Expected 1 shrink, got %d", len(shrinks))
	}

	// Should have second level shrinks
	secondLevel := shrinks[0].Shrinks()
	if len(secondLevel) != 1 {
		t.Errorf("Expected 1 second level shrink, got %d", len(secondLevel))
	}

	// Should NOT have third level shrinks
	thirdLevel := secondLevel[0].Shrinks()
	if len(thirdLevel) != 0 {
		t.Errorf("Expected 0 third level shrinks (pruned), got %d", len(thirdLevel))
	}
}

func TestTreeFindSmallest(t *testing.T) {
	tree := NewTree(10, func() []*Tree[int] {
		return []*Tree[int]{
			NewTree(5, func() []*Tree[int] {
				return []*Tree[int]{
					NewTree(2, func() []*Tree[int] { return nil }),
				}
			}),
			NewTree(8, func() []*Tree[int] { return nil }),
		}
	})

	// Find smallest value that's > 3
	smallest, found := tree.FindSmallest(func(x int) bool { return x > 3 })

	if !found {
		t.Error("Expected to find a value > 3")
	}

	if smallest != 5 {
		t.Errorf("Expected smallest value > 3 to be 5, got %d", smallest)
	}

	// Find smallest value that's > 10 (should not exist)
	_, found = tree.FindSmallest(func(x int) bool { return x > 10 })
	if found {
		t.Error("Expected not to find a value > 10")
	}
}

// Edge cases and robustness tests

func TestTreeUnfoldCycleDetection(t *testing.T) {
	// Create a tree that could potentially create cycles
	tree := NewTree(1, func() []*Tree[int] {
		return []*Tree[int]{
			NewTree(2, func() []*Tree[int] {
				return []*Tree[int]{
					NewTree(3, func() []*Tree[int] {
						// This creates a very deep tree
						return []*Tree[int]{
							NewTree(4, func() []*Tree[int] { return nil }),
						}
					}),
				}
			}),
		}
	})

	// Should not hang or crash
	result := tree.Unfold()

	// Should have limited results due to safety bounds
	if len(result) == 0 {
		t.Error("Expected some results from Unfold")
	}

	if len(result) > 10000 {
		t.Error("Unfold should limit results to prevent memory issues")
	}
}

func TestTreeFilterAllRejected(t *testing.T) {
	tree := NewTree(10, func() []*Tree[int] {
		return []*Tree[int]{
			NewTree(5, func() []*Tree[int] { return nil }),
			NewTree(3, func() []*Tree[int] { return nil }),
		}
	})

	// Filter that rejects everything
	filtered := tree.Filter(func(x int) bool { return false })

	// Root should still be there
	if filtered.Value() != 10 {
		t.Errorf("Expected filtered value 10, got %d", filtered.Value())
	}

	// But no shrinks should remain
	shrinks := filtered.Shrinks()
	if len(shrinks) != 0 {
		t.Errorf("Expected no shrinks after filtering all, got %d", len(shrinks))
	}
}

func TestTreePruneZeroDepth(t *testing.T) {
	tree := NewTree(10, func() []*Tree[int] {
		return []*Tree[int]{
			NewTree(5, func() []*Tree[int] { return nil }),
		}
	})

	// Prune to zero depth
	pruned := tree.Prune(0)

	if pruned.Value() != 10 {
		t.Errorf("Expected pruned value 10, got %d", pruned.Value())
	}

	// Should have no shrinks
	shrinks := pruned.Shrinks()
	if len(shrinks) != 0 {
		t.Errorf("Expected no shrinks after pruning to depth 0, got %d", len(shrinks))
	}
}
