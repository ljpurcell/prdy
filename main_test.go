package main

import (
	"fmt"
	"testing"
)

func TestSliceElementIsSubstringWhenTrue(t *testing.T) {
	sliceOfSubstrings := []string{"this", "is", "my", "slice"}
	stringToSearch := "slices"

	expected := true
	actual := SliceElementIsSubstring(sliceOfSubstrings, stringToSearch)

	if actual != expected {
		t.Errorf("Expected: %v. Actual: %v. Using slice of %v and string of %q.", expected, actual, sliceOfSubstrings, stringToSearch)
	}
}

func TestSliceElementIsSubstringWhenFalse(t *testing.T) {
	sliceOfSubstrings := []string{"this", "is", "my", "slice"}
	stringToSearch := "bananas"

	expected := false
	actual := SliceElementIsSubstring(sliceOfSubstrings, stringToSearch)

	if actual != expected {
		t.Errorf("Expected: %v. Actual: %v. Using slice of %v and string of %q.", expected, actual, sliceOfSubstrings, stringToSearch)
	}
}

func TestSearchConfigAddMethod(t *testing.T) {
	sc := SearchConfig{[]string{"var_dump", "dd", "console.log"}, []string{}, []string{}, []string{}}
	before := len(sc.HitWords)
	newElements := []string{"rand", "new element"}
	sc.addToField(newElements, &sc.HitWords)
	expected := before + len(newElements)
	actual := len(sc.HitWords)

	if actual != expected {
		t.Errorf("Expected: %v. Actual: %v. Length of %v after adding %v.", expected, actual, sc.HitWords, newElements)
	}
}

func TestStringToIntSliceSorted(t *testing.T) {
	expected := []int{1, 2, 3, 4, 5}
	actual := ToIntSorted([]string{"1", "5", "4", "2", "3"})

	if len(expected) != len(actual) {
		t.Errorf("Length - Expected: %v. Actual: %v.", expected, actual)
	}

	for i := range expected {
		if actual[i] != expected[i] {
			t.Errorf("Expected %q but got %q at index %d", expected[i], actual[i], i)
		}
	}
}

func TestSearchConfigRemoveMethod(t *testing.T) {
	sc := SearchConfig{[]string{"var_dump", "dd", "console.log"}, []string{}, []string{}, []string{}}
	before := len(sc.HitWords)
	elementIndexToRemove := []string{"1", "2"}
	removedElements, err := sc.removeFromField(elementIndexToRemove, &sc.HitWords)
	if err != nil {
		fmt.Printf("Error in TestSearchConfigRemoveMethod. %s\n", err)
	}
	expected := before - len(removedElements)
	actual := len(sc.HitWords)

	if actual != expected {
		t.Errorf("Expected: %v. Actual: %v. Length of %v after removing %q.", expected, actual, sc.HitWords, removedElements)
	}

	for i, v := range sc.HitWords {
		for _, removed := range removedElements {
			if v == removed {
				t.Errorf("Found %q at index %d. Should have been removed.", removed, i)
			}
		}
	}
}
