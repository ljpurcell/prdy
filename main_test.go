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

func TestSearchConfigAddMethod(t *testing.T) {
	sc := SearchConfig{[]string{"var_dump", "dd", "console.log"}, []string{}, []string{}, []string{}}
	before := len(sc.HitWords)
	newElement := "rand"
	sc.addToField(newElement, &sc.HitWords)
	expected := before + 1
	actual := len(sc.HitWords)

	if actual != expected {
		t.Errorf("Expected: %v. Actual: %v. Length of %v after adding %q.", expected, actual, sc.HitWords, newElement)
	}
}

func TestSearchConfigRemoveMethod(t *testing.T) {
	sc := SearchConfig{[]string{"var_dump", "dd", "console.log"}, []string{}, []string{}, []string{}}
	before := len(sc.HitWords)
	elementIndexToRemove := 1
	removedElement := sc.removeFromField(elementIndexToRemove, &sc.HitWords)
	expected := before - 1
	actual := len(sc.HitWords)

	if actual != expected {
		t.Errorf("Expected: %v. Actual: %v. Length of %v after removing %q.", expected, actual, sc.HitWords, removedElement)
	}

	for i, v := range sc.HitWords {
		if v == removedElement {
			t.Errorf("Found %q at index %d. Should have been removed.", removedElement, i)
		}
	}
}
