package main

import (
	"io/fs"
	"sort"
	"strconv"
	"strings"
)

func ToIntSorted(input []string) []int {
	var result []int
	for _, strV := range input {
		intV, err := strconv.ParseInt(strV, 10, strconv.IntSize)
		Check(err)
		result = append(result, int(intV))
	}
	sort.Ints(result)
	return result
}

/*
 * ERROR HANDLING
 * TODO: consider creating a varienty of more specific error handlers
 * Do some research on what more specific EH would look like
 */
func Check(err error) {
	if err != nil {
		panic(err)
	}
}

func SliceElementIsSubstring(needles []string, haystack string) bool {
	for _, needle := range needles {
		if strings.Contains(haystack, needle) {
			return true
		}
	}
	return false
}

/*
 * FILE HANDLING
 */
func CloseFile(file fs.File) {
	err := file.Close()
	Check(err)
}

/*
 * GLOB MATCHING
 */
// https://github.com/bmatcuk/doublestar
