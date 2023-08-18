package main

import (
	"bufio"
	"fmt"
	"io/fs"
)

/*
 * PRIMARY FUNCTIONALITY
 */
func CheckFileForHits(file fs.File, sc *SearchConfig) {
	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	var outputMap = make(map[string][]string)
	var outputArray []string

	fileInfo, err := file.Stat()
	Check(err)
	fileName := fileInfo.Name()

	for i := 1; fileScanner.Scan(); i++ {
		// Need to find a single hitword & find zero excluded versions
		if foundMatch(sc.HitWords, fileScanner.Text()) && !foundMatch(sc.ExcludedWords, fileScanner.Text()) {
			line := fmt.Sprintf("%v %v\n", i, fileScanner.Text())
			outputArray = append(outputArray, line)
		}
	}

	outputMap[fileName] = outputArray

	for i, line := range outputArray {
		if i == 0 {
			fmt.Printf("\n\tFile: %s", fileName)
		}
		fmt.Printf(line)
	}

}
