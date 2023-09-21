package main

import (
	"bufio"
	"fmt"
	"io/fs"
)

/*
 * PRIMARY FUNCTIONALITY
 */
func CheckFileForHits(file fs.File, m map[string][]string, sc *SearchConfig) map[string][]string {

	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	var hitLines []string

	fileInfo, err := file.Stat()
	Check(err)
	fileName := fileInfo.Name()

	for i := 1; fileScanner.Scan(); i++ {
		// Need to find a single hitword & find zero excluded versions
		if foundMatch(sc.HitWords, fileScanner.Text()) && !foundMatch(sc.ExcludedWords, fileScanner.Text()) {
			line := fmt.Sprintf("line %v: %v\n", i, fileScanner.Text())
			hitLines = append(hitLines, line)
		}
	}

	m[fileName] = hitLines

	return m
}

func DisplayHitsForEachFile(m map[string][]string) {

	for file := range m {

		fmt.Printf("\n\tFile: %v\n", file)

		for line := range m[file] {

			fmt.Printf("%v\n", line)
		}
	}
}
