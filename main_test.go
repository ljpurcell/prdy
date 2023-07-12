package main

import (
	"testing"
)

func TestSliceElementIsSubstringWhenTrue(t *testing.T) {
	sliceOfSubstrings := []string{"this", "is", "my", "slice"}
	stringToSearch := "slices"

	expected := true
	actual := slice_element_is_substring(sliceOfSubstrings, stringToSearch)

	if actual != expected {
		t.Errorf("Expected: %v. Actual: %v. Using slice of %v and string of %q.", expected, actual, sliceOfSubstrings, stringToSearch)
	}
}

func TestSliceElementIsSubstringWhenFalse(t *testing.T) {
	sliceOfSubstrings := []string{"this", "is", "my", "slice"}
	stringToSearch := "bananas"

	expected := false
	actual := slice_element_is_substring(sliceOfSubstrings, stringToSearch)

	if actual != expected {
		t.Errorf("Expected: %v. Actual: %v. Using slice of %v and string of %q.", expected, actual, sliceOfSubstrings, stringToSearch)
	}

}
