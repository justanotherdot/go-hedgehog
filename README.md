# Go Hedgehog

> Release with confidence.

Property-based testing library for Go, inspired by the original [Hedgehog](https://hedgehog.qa/) library for Haskell.

## Features

- **Explicit generators** - No type-directed magic, generators are first-class values you compose
- **Integrated shrinking** - Shrinks obey invariants by construction, built into generators
- **Excellent debugging** - Minimal counterexamples with rich failure reporting
- **Distribution shaping** - Control probability distributions for realistic test data
- **Variable name tracking** - Enhanced failure reporting with named inputs
- **Property classification** - Inspect test data distribution and statistics

## Quick Start

```go
package main

import (
    "testing"
    "github.com/rjs/go-hedgehog"
)

func TestReverse(t *testing.T) {
    gen := hedgehog.SliceOf(hedgehog.IntRange(1, 100))
    prop := hedgehog.ForAll(gen, func(xs []int) bool {
        reversed := reverse(xs)
        doubleReversed := reverse(reversed)
        return slicesEqual(xs, doubleReversed)
    })
    
    if !prop.Run(hedgehog.DefaultConfig()) {
        t.Error("Property failed")
    }
}
```

## Project Status

This is a work-in-progress implementation, porting the Rust hedgehog library to Go.

## License

This project is licensed under the BSD-3-Clause License.

## Acknowledgments

- Jacob Stanley and the original Hedgehog team for the foundational ideas
- The Haskell, F#, Rust, and R Hedgehog ports for implementation insights