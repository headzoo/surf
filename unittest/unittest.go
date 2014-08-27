package unittest

import (
	"reflect"
	"strings"
	"testing"
)

// Test is the testing.T instance for the test being run.
var test *testing.T = &testing.T{}

// Run sets up an individual test.
func Run(t *testing.T) {
	test = t
}

// Test returns the testing.T passed to the Run() method.
func Test() *testing.T {
	return test
}

// AssertTrue tests whether the given value is true.
func AssertTrue(actual bool) bool {
	if !actual {
		test.Errorf("Failed asserting %q is true.", actual)
		return false
	}
	return true
}

// AssertFalse tests whether the given value is false.
func AssertFalse(actual bool) bool {
	if actual {
		test.Errorf("Failed asserting %q is false.", actual)
		return false
	}
	return true
}

// AssertNotNil tests whether the given value is not nil.
func AssertNotNil(actual interface{}) bool {
	if actual == nil {
		test.Errorf("Failed asserting the value is not nil.")
		return false
	}
	return true
}

// AssertNil tests whether the given value is nil.
func AssertNil(actual interface{}) bool {
	if actual != nil {
		test.Errorf("Failed asserting %T is nil.", actual)
		return false
	}
	return true
}

// AssertEmpty tests whether the given value is empty.
func AssertEmpty(actual interface{}) bool {
	t := reflect.ValueOf(actual)
	if !isZero(t) {
		test.Errorf("Failed asserting %q is empty.", actual)
		return false
	}
	return true
}

// AssertNotEmpty tests whether the given value is not empty.
func AssertNotEmpty(actual interface{}) bool {
	t := reflect.ValueOf(actual)
	if isZero(t) {
		test.Errorf("Failed asserting %q is not empty.", actual)
		return false
	}
	return true
}

// AssertEquals tests whether two values are equal.
func AssertEquals(expected, actual interface{}) bool {
	if !reflect.DeepEqual(expected, actual) {
		test.Errorf("Failed asserting %q equals %q.", expected, actual)
		return false
	}
	return true
}

// AssertNotEquals tests whether two values do not equal each other.
func AssertNotEquals(expected, actual interface{}) bool {
	if reflect.DeepEqual(expected, actual) {
		test.Errorf("Failed asserting %q is not equal to %q.", expected, actual)
		return false
	}
	return true
}

// AssertGreaterThan tests whether the actual value is greater than the expected value.
func AssertGreaterThan(expected, actual int) bool {
	if expected >= actual {
		test.Errorf("Failed asserting %q is greater than %q.", actual, expected)
		return false
	}
	return true
}

// AssertContains tests whether the expected value contains the actual value.
func AssertContains(expected, actual string) bool {
	if !strings.Contains(actual, expected) {
		test.Errorf("Failed asserting %q contains %q.", actual, expected)
		return false
	}
	return true
}

// isZero returns if the value is zero.
func isZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Func, reflect.Map, reflect.Slice:
		return v.IsNil()
	case reflect.Array:
		z := true
		for i := 0; i < v.Len(); i++ {
			z = z && isZero(v.Index(i))
		}
		return z
	case reflect.Struct:
		z := true
		for i := 0; i < v.NumField(); i++ {
			z = z && isZero(v.Field(i))
		}
		return z
	}
	// Compare other types directly:
	z := reflect.Zero(v.Type())
	return v.Interface() == z.Interface()
}
