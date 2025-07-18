package hedgehog

import (
	"fmt"
	"strings"
	"time"
)

// Property represents a property that can be tested.
type Property[T any] struct {
	generator *Gen[T]
	predicate func(T) bool
	name      string
}

// NewProperty creates a new property with the given generator and predicate.
func NewProperty[T any](generator *Gen[T], predicate func(T) bool) *Property[T] {
	return &Property[T]{
		generator: generator,
		predicate: predicate,
		name:      "",
	}
}

// ForAll creates a property that tests the predicate for all generated values.
func ForAll[T any](generator *Gen[T], predicate func(T) bool) *Property[T] {
	return NewProperty(generator, predicate)
}

// ForAllNamed creates a property with a name for better error reporting.
func ForAllNamed[T any](generator *Gen[T], name string, predicate func(T) bool) *Property[T] {
	prop := NewProperty(generator, predicate)
	prop.name = name
	return prop
}

// WithName sets the name of the property.
func (p *Property[T]) WithName(name string) *Property[T] {
	p.name = name
	return p
}

// Run executes the property test with the given configuration.
func (p *Property[T]) Run(config *Config) *TestReport[T] {
	return p.RunWithReport(config)
}

// RunWithReport executes the property test and returns a detailed report.
func (p *Property[T]) RunWithReport(config *Config) *TestReport[T] {
	random := NewRandom(config.Seed)

	report := &TestReport[T]{
		PropertyName: p.name,
		Config:       config,
		StartTime:    time.Now(),
	}

	for i := 0; i < config.TestCount; i++ {
		tree := p.generator.Generate(random)

		if p.predicate(tree.Value()) {
			report.PassedTests++
			continue
		}

		// Property failed, try to find minimal counterexample
		counterexample, shrinkSteps := p.findMinimalCounterexample(tree, config)

		report.Result = TestFail
		report.Counterexample = counterexample
		report.ShrinkSteps = shrinkSteps
		report.FailedAfter = i + 1
		report.EndTime = time.Now()

		return report
	}

	// All tests passed
	report.Result = TestPass
	report.EndTime = time.Now()
	return report
}

// findMinimalCounterexample finds the smallest counterexample through shrinking.
func (p *Property[T]) findMinimalCounterexample(tree *Tree[T], config *Config) (T, []T) {
	var shrinkSteps []T
	current := tree

	// Keep shrinking until we can't shrink any further
	for {
		shrinks := current.Shrinks()
		foundSmaller := false

		for _, shrink := range shrinks {
			if !p.predicate(shrink.Value()) {
				shrinkSteps = append(shrinkSteps, shrink.Value())
				current = shrink
				foundSmaller = true
				break
			}
		}

		if !foundSmaller {
			break
		}

		// Prevent infinite shrinking
		if len(shrinkSteps) >= config.ShrinkLimit {
			break
		}
	}

	return current.Value(), shrinkSteps
}

// TestReport contains the results of running a property test.
type TestReport[T any] struct {
	PropertyName   string
	Config         *Config
	Result         TestResult
	PassedTests    int
	FailedAfter    int
	Counterexample T
	ShrinkSteps    []T
	StartTime      time.Time
	EndTime        time.Time
}

// Duration returns the time taken to run the test.
func (r *TestReport[T]) Duration() time.Duration {
	return r.EndTime.Sub(r.StartTime)
}

// Success returns true if the test passed.
func (r *TestReport[T]) Success() bool {
	return r.Result == TestPass
}

// String returns a string representation of the test report.
func (r *TestReport[T]) String() string {
	var sb strings.Builder

	switch r.Result {
	case TestPass:
		sb.WriteString(fmt.Sprintf("✓ Property passed %d tests", r.PassedTests))
		if r.PropertyName != "" {
			sb.WriteString(fmt.Sprintf(" (%s)", r.PropertyName))
		}
		sb.WriteString(fmt.Sprintf(" in %v\n", r.Duration()))

	case TestFail:
		sb.WriteString(fmt.Sprintf("✗ Property failed after %d tests", r.FailedAfter))
		if r.PropertyName != "" {
			sb.WriteString(fmt.Sprintf(" (%s)", r.PropertyName))
		}
		sb.WriteString(fmt.Sprintf(" in %v\n", r.Duration()))

		if len(r.ShrinkSteps) > 0 {
			sb.WriteString(fmt.Sprintf("  Shrinking progression (%d steps):\n", len(r.ShrinkSteps)))
			for i, step := range r.ShrinkSteps {
				sb.WriteString(fmt.Sprintf("    %d: %v\n", i+1, step))
			}
		}

		sb.WriteString(fmt.Sprintf("  Minimal counterexample: %v\n", r.Counterexample))

	case TestGaveUp:
		sb.WriteString("? Property gave up (too many discarded cases)\n")
	}

	return sb.String()
}

// Print prints the test report to stdout.
func (r *TestReport[T]) Print() {
	fmt.Print(r.String())
}

// Assert fails the test if the property didn't pass.
func (r *TestReport[T]) Assert(t TestingT) {
	if !r.Success() {
		t.Error(r.String())
	}
}

// TestingT is an interface that matches testing.T for integration with Go's testing package.
type TestingT interface {
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Helper()
}
