package hedgehog

import (
	"testing"
)

func TestPropertyRun(t *testing.T) {
	gen := IntRange(1, 10)
	prop := ForAll(gen, func(x int) bool {
		return x >= 1 && x <= 10
	})
	
	config := &Config{
		TestCount:   50,
		ShrinkCount: 10,
		Seed:        12345,
	}
	
	report := prop.Run(config)
	
	if !report.Success() {
		t.Error("Property should have passed")
	}
	
	if report.PassedTests != 50 {
		t.Errorf("Expected 50 passed tests, got %d", report.PassedTests)
	}
}

func TestPropertyFail(t *testing.T) {
	gen := IntRange(1, 10)
	prop := ForAll(gen, func(x int) bool {
		return x < 5 // This will fail for values >= 5
	})
	
	config := &Config{
		TestCount:   100,
		ShrinkCount: 10,
		Seed:        12345,
	}
	
	report := prop.Run(config)
	
	if report.Success() {
		t.Error("Property should have failed")
	}
	
	if report.Result != TestFail {
		t.Errorf("Expected TestFail, got %v", report.Result)
	}
	
	if report.Counterexample < 5 {
		t.Errorf("Expected counterexample >= 5, got %d", report.Counterexample)
	}
	
	if report.FailedAfter == 0 {
		t.Error("Expected FailedAfter > 0")
	}
}

func TestPropertyShrinking(t *testing.T) {
	gen := SliceOf(IntRange(1, 100))
	prop := ForAll(gen, func(xs []int) bool {
		return len(xs) < 3 // This will fail for slices with 3+ elements
	})
	
	config := &Config{
		TestCount:   100,
		ShrinkCount: 100,
		Seed:        12345,
	}
	
	report := prop.Run(config)
	
	if report.Success() {
		t.Error("Property should have failed")
	}
	
	// Should have found a counterexample
	if len(report.Counterexample) < 3 {
		t.Errorf("Expected counterexample with >= 3 elements, got %d", len(report.Counterexample))
	}
	
	// Should have shrunk to a smaller example
	if len(report.ShrinkSteps) == 0 {
		t.Error("Expected some shrink steps")
	}
	
	// Final counterexample should be minimal (exactly 3 elements)
	if len(report.Counterexample) != 3 {
		t.Errorf("Expected minimal counterexample with 3 elements, got %d", len(report.Counterexample))
	}
}

func TestPropertyNamed(t *testing.T) {
	gen := IntRange(1, 10)
	prop := ForAllNamed(gen, "value", func(x int) bool {
		return x >= 1 && x <= 10
	})
	
	config := DefaultConfig()
	config.TestCount = 10
	
	report := prop.Run(config)
	
	if !report.Success() {
		t.Error("Named property should have passed")
	}
	
	if report.PropertyName != "value" {
		t.Errorf("Expected property name 'value', got %q", report.PropertyName)
	}
}

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()
	
	if config.TestCount != 100 {
		t.Errorf("Expected TestCount 100, got %d", config.TestCount)
	}
	
	if config.ShrinkCount != 100 {
		t.Errorf("Expected ShrinkCount 100, got %d", config.ShrinkCount)
	}
	
	if config.Seed == 0 {
		t.Error("Expected non-zero seed")
	}
}

func TestReportString(t *testing.T) {
	gen := IntRange(1, 10)
	prop := ForAllNamed(gen, "x", func(x int) bool {
		return x < 5
	})
	
	config := &Config{
		TestCount:   10,
		ShrinkCount: 10,
		Seed:        12345,
	}
	
	report := prop.Run(config)
	
	if report.Success() {
		// Test pass report
		str := report.String()
		if str == "" {
			t.Error("Expected non-empty string representation")
		}
	} else {
		// Test fail report
		str := report.String()
		if str == "" {
			t.Error("Expected non-empty string representation")
		}
		
		// Should contain counterexample info
		if len(report.ShrinkSteps) > 0 && !contains(str, "Shrinking progression") {
			t.Error("Expected shrinking progression in failure report")
		}
		
		if !contains(str, "Minimal counterexample") {
			t.Error("Expected minimal counterexample in failure report")
		}
	}
}

func TestReportDuration(t *testing.T) {
	gen := IntRange(1, 10)
	prop := ForAll(gen, func(x int) bool {
		return x >= 1 && x <= 10
	})
	
	config := &Config{
		TestCount:   5,
		ShrinkCount: 5,
		Seed:        12345,
	}
	
	report := prop.Run(config)
	
	duration := report.Duration()
	if duration <= 0 {
		t.Error("Expected positive duration")
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[len(s)-len(substr):] == substr || 
		   len(s) > len(substr) && s[:len(substr)] == substr ||
		   containsMiddle(s, substr)
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}