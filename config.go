package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

/**
 * CONFIGUATION STRUCT
 */
type SearchConfig struct {
	HitWords          []string
	ExcludedWords     []string
	SourceDirectories []string
	IgnoredFiles      []string
}

/**
 * METHODS
 */
func (sc *SearchConfig) addToField(elements []string, field *[]string) {
	defer sc.updateConfigFile()
	for _, element := range elements {
		*field = append(*field, element)
	}
}

func (sc *SearchConfig) removeFromField(indicesStr []string, field *[]string) ([]string, error) {
	defer sc.updateConfigFile()
	var removedElements []string
	var err error

	indices := ToIntSorted(indicesStr)

	for i, index := range indices {
		toRemove := index - i
		removedElements = append(removedElements, (*field)[toRemove])
		copy((*field)[toRemove:], (*field)[toRemove+1:])
		(*field) = (*field)[:len((*field))-1]
	}
	return removedElements, err
}

func (sc *SearchConfig) updateConfigFile() {
	scJson, err := json.Marshal(*sc)
	Check(err)
	os.WriteFile(".prdy_config.json", scJson, 0644)
}

func removeExcludedWord(sc *SearchConfig) {
	fmt.Println("\n\t* Remove Hit Word *")
	fmt.Println("\nPlease type the number of the word you want to remove.")
	fmt.Println("TIP: If you want to remove mutliple words, type a space seperated list.")

	displayExcludedVersions(sc, true)

	fmt.Print("\nRemove excluded word: ")
	var indicesToRemove string
	inputReader := bufio.NewReader(os.Stdin)
	indicesToRemove, err := inputReader.ReadString('\n')
	Check(err)

	providedIndices := strings.Split(indicesToRemove, " ")

	for i, indexString := range providedIndices {
		indexString = strings.TrimSpace(indexString)
		indexValue, err := strconv.ParseInt(indexString, 10, strconv.IntSize)
		Check(err)

		indexValue -= 1 + int64(i) // because menu is 1-based and indices become progressively off by one more each time around the loop as an item is removed
		removedWord := sc.ExcludedWords[indexValue]
		copy(sc.ExcludedWords[indexValue:], sc.ExcludedWords[indexValue+1:])
		sc.ExcludedWords = sc.ExcludedWords[:len(sc.ExcludedWords)-1]

		fmt.Printf("Removed %q\n", removedWord)

	}
}

func removeIgnoredFile(sc *SearchConfig) {

	// defer sc.updateConfigFile(sc)

	fmt.Println("\n\t* Remove Ignored File *")
	fmt.Println("\nPlease type the number of the word you want to remove.")
	fmt.Println("TIP: If you want to remove multiple files, type a space seperated list.")

	displayIgnoredFiles(sc, true)

	fmt.Print("\nRemove ignored files: ")
	var indicesToRemove string
	inputReader := bufio.NewReader(os.Stdin)
	indicesToRemove, err := inputReader.ReadString('\n')
	Check(err)

	providedIndices := strings.Split(indicesToRemove, " ")

	for i, indexString := range providedIndices {
		indexString = strings.TrimSpace(indexString)
		indexValue, err := strconv.ParseInt(indexString, 10, strconv.IntSize)
		Check(err)

		indexValue -= 1 + int64(i) // because menu is 1-based and indices become progressively off by one more each time around the loop as an item is removed
		removedWord := sc.IgnoredFiles[indexValue]
		copy(sc.IgnoredFiles[indexValue:], sc.IgnoredFiles[indexValue+1:])
		sc.IgnoredFiles = sc.IgnoredFiles[:len(sc.IgnoredFiles)-1]

		fmt.Printf("Removed %q\n", removedWord)
	}
}

func getGitIgnorePatterns() []string {
	file, err := os.Open(".gitignore")
	Check(err)
	defer CloseFile(file)

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var ignoredPatterns []string
	for scanner.Scan() {

		// comments in gitignore file
		if !strings.HasPrefix(scanner.Text(), "#") && len(strings.TrimSpace(scanner.Text())) != 0 {
			ignoredPatterns = append(ignoredPatterns, scanner.Text())
		}
	}
	return ignoredPatterns
}

func loadConfig() *SearchConfig {
	configJson, err := os.ReadFile(".prdy_config.json")
	Check(err)
	var sc SearchConfig
	json.Unmarshal(configJson, &sc)
	return &sc
}

func createEmptyConfig() *SearchConfig {
	sc := &SearchConfig{}
	jsonSc, err := json.Marshal(sc)
	Check(err)
	os.WriteFile(".prdy_config.json", jsonSc, 0644)
	return sc
}
