# Go Hedgehog roadmap

This document outlines the development plan for go-hedgehog, tracking features from the Rust hedgehog implementation that need to be ported.

## Current status

The Go port currently implements ~30% of the full Rust hedgehog feature set, covering the essential property testing core with basic generators and shrinking.

### Completed features

- Basic generators: `IntRange`, `Float64Range`, `Bool`, `String`, `SliceOf`
- Generator combinators: `OneOf`, `Frequency`, `Map`
- Integrated shrinking: Automatic minimal counterexample finding
- Rich reporting: Detailed failure reports with shrinking progression
- Testing integration: Works with standard Go testing framework
- Lazy shrinking trees with proper evaluation

## Phase 1: Core features (high priority)

### Range/distribution system
- [ ] `Range` type with uniform, linear, exponential distributions
- [ ] `Gen.FromRange()` for distribution-controlled generation
- [ ] Probability distribution shaping for realistic test data

### Variable name tracking
- [ ] Enhanced `ForAllNamed()` with proper variable name display
- [ ] Format failure reports with `-- variable_name` syntax
- [ ] Support for multiple named variables

### Property classification
- [ ] `.Classify()` method for categorizing test data
- [ ] `.Collect()` method for gathering statistics
- [ ] Statistics reporting (min, max, avg, median)
- [ ] Test data distribution analysis

### Enhanced error types
- [ ] Structured error types following Go conventions
- [ ] External vs internal error differentiation
- [ ] Recoverable vs irrecoverable error handling

## Phase 2: Generator expansion (medium priority)

### Advanced string generators
- [ ] Character set generators: `Alpha()`, `ASCII()`, `Unicode()`
- [ ] Length distribution control
- [ ] Custom character set support
- [ ] String pattern generators

### Numeric generators
- [ ] Type-specific generators: `Int8()`, `Int16()`, `Int32()`, `Int64()`
- [ ] Unsigned integer generators: `Uint8()`, `Uint16()`, `Uint32()`, `Uint64()`
- [ ] Float precision generators: `Float32()`, `Float64()`
- [ ] Size and bounds configuration

### Collection generators
- [ ] `MapOf()` generator for `map[K]V` types
- [ ] `TupleOf()` generators for multiple element tuples
- [ ] `StructOf()` for anonymous struct generation
- [ ] Nested collection support

## Phase 3: Advanced testing paradigms (medium priority)

### State machine testing
- [ ] State machine testing utilities
- [ ] Test stateful systems systematically
- [ ] State transition validation
- [ ] Action sequence generation

### Function generators
- [ ] Generate functions as test inputs
- [ ] Higher-order property testing
- [ ] Function behavior validation
- [ ] Lambda/closure generation

### Coverage-guided generation
- [ ] Use coverage feedback to guide test generation
- [ ] Integrate with Go's coverage tooling
- [ ] Adaptive test case generation
- [ ] Coverage-directed shrinking

### Example integration
- [ ] Mix explicit examples with generated tests
- [ ] Hybrid testing approach
- [ ] Example-driven generation
- [ ] Regression case integration

### Dictionary support
- [ ] Domain-specific token injection
- [ ] Realistic data generation for specific domains
- [ ] Custom vocabulary support
- [ ] Context-aware generation

## Phase 4: Production features (lower priority)

### Parallel testing
- [ ] Find race conditions through parallel property testing
- [ ] Concurrent test execution
- [ ] Goroutine-based testing
- [ ] Data race detection

### Fault injection
- [ ] Systematic failure testing
- [ ] Inject failures at specific points
- [ ] Error scenario generation
- [ ] Resilience testing

### Reflection-based generators
- [ ] Automatic generator creation for structs using reflection
- [ ] Interface-based generator selection
- [ ] Custom generator registration system

### Corpus/regression testing
- [ ] Test corpus persistence to disk
- [ ] Regression test functionality
- [ ] Corpus minimization and management

### Enhanced configuration
- [ ] Verbose output modes
- [ ] Custom reporting formats
- [ ] Performance profiling options
- [ ] Distributed test execution

### Tree rendering improvements
- [ ] Compact rendering mode
- [ ] Numbered shrink progression
- [ ] Shrink-specific tree views
- [ ] Graphical tree visualization

## Phase 5: Go-specific enhancements

### Integration improvements
- [ ] Benchmark integration with `testing.B`
- [ ] Fuzzing integration with Go 1.18+ fuzzing
- [ ] Custom assertion helpers
- [ ] Test table integration

### Performance optimizations
- [ ] Memory-efficient shrinking
- [ ] Parallel property evaluation
- [ ] Streaming test execution
- [ ] Benchmark-driven optimizations

### Tooling
- [ ] CLI tool for corpus management
- [ ] IDE integration helpers
- [ ] Test coverage analysis
- [ ] Performance monitoring

## Implementation notes

### Go-specific considerations
- **Generics**: Use Go 1.18+ generics for type safety
- **Reflection**: Leverage `reflect` package for automatic generator creation
- **Interfaces**: Design clean interfaces for extensibility
- **Error handling**: Follow Go error conventions with structured error types
- **Testing**: Deep integration with Go's testing framework

### Architecture principles
- **Composability**: Generators should be easily composable
- **Laziness**: Shrinking trees should be lazy by default
- **Performance**: Optimize for both speed and memory usage
- **Simplicity**: Keep the API simple and discoverable

## Contributing

Features are roughly prioritized by impact and implementation difficulty. Phase 1 items are essential for a complete property testing library, while later phases add convenience and advanced features.

Each feature should include:
- Implementation in the core library
- Comprehensive tests
- Documentation and examples
- Performance benchmarks where applicable

## Future considerations

- **WebAssembly support**: Ensure compatibility with WASM targets
- **Cross-platform testing**: Verify behavior across different Go versions and platforms
- **Ecosystem integration**: Consider integration with popular Go testing libraries
- **Language server support**: Provide better IDE integration for generator development