package config

import (
    "fmt"
    "github.com/ljpurcell/prdy/utils"

)

/**
 * DATA TYPE
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

	indices := toIntSorted(indicesStr)

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
	check(err)
	os.WriteFile(".prdy_config.json", scJson, 0644)
}

func addHitWords(sc *SearchConfig) {
	defer sc.updateConfigFile()

	fmt.Println("\n\t* Add Hit Words *")
	fmt.Println("\nPlease type the pattern you want to match on. To add multiple, use a space seperated list.")
	fmt.Println("TIP: If you are looking for a function, leave off the parenthesis -- unless you know the exact naming of the argument(s) it has been called with.")
	fmt.Print("Add hit word: ")
	var input string
	fmt.Scanln(&input)
	wordsToAdd := strings.Split(input, " ")
	sc.addToField(wordsToAdd, &sc.HitWords)
}

func removeHitWord(sc *SearchConfig) {
	// defer sc.updateConfigFile(sc)

	fmt.Println("\n\t* Remove Hit Word *")
	fmt.Println("\nPlease type the number of the word you want to remove.")
	fmt.Println("TIP: If you want to remove mutliple words, type a space seperated list.")

	displayHitWords(sc, true)

	fmt.Print("\nRemove hit word: ")
	var indicesToRemove string
	inputReader := bufio.NewReader(os.Stdin)
	indicesToRemove, err := inputReader.ReadString('\n')
	check(err)

	providedIndices := strings.Split(indicesToRemove, " ")

	for i, indexString := range providedIndices {
		indexString = strings.TrimSpace(indexString)
		indexValue, err := strconv.ParseInt(indexString, 10, strconv.IntSize)
		check(err)

		indexValue -= 1 + int64(i) // because menu is 1-based and indices become progressively off by one more each time around the loop as an item is removed
		removedWord := sc.HitWords[indexValue]
		copy(sc.HitWords[indexValue:], sc.HitWords[indexValue+1:])
		sc.HitWords = sc.HitWords[:len(sc.HitWords)-1]

		fmt.Printf("Removed %q\n", removedWord)
	}
}

func addExcludedWords(sc *SearchConfig) {
	// defer sc.updateConfigFile(sc)

	fmt.Println("\n\t* Adding new Excluded Word *")
	fmt.Println("\nPlease type the word or pattern you want to PREVENT matching on.")
	fmt.Println("TIP: This is where you can use specific argument names to stop returning false positives.")
	fmt.Print("Add excluded word: ")
	var input string
	fmt.Scanln(&input)
	wordsToAdd := strings.Split(input, " ")
	sc.addToField(wordsToAdd, &sc.ExcludedWords)
}

func removeExcludedWord(sc *SearchConfig) {
	// defer sc.updateConfigFile(sc)

	fmt.Println("\n\t* Remove Hit Word *")
	fmt.Println("\nPlease type the number of the word you want to remove.")
	fmt.Println("TIP: If you want to remove mutliple words, type a space seperated list.")

	displayExcludedVersions(sc, true)

	fmt.Print("\nRemove excluded word: ")
	var indicesToRemove string
	inputReader := bufio.NewReader(os.Stdin)
	indicesToRemove, err := inputReader.ReadString('\n')
	check(err)

	providedIndices := strings.Split(indicesToRemove, " ")

	for i, indexString := range providedIndices {
		indexString = strings.TrimSpace(indexString)
		indexValue, err := strconv.ParseInt(indexString, 10, strconv.IntSize)
		check(err)

		indexValue -= 1 + int64(i) // because menu is 1-based and indices become progressively off by one more each time around the loop as an item is removed
		removedWord := sc.ExcludedWords[indexValue]
		copy(sc.ExcludedWords[indexValue:], sc.ExcludedWords[indexValue+1:])
		sc.ExcludedWords = sc.ExcludedWords[:len(sc.ExcludedWords)-1]

		fmt.Printf("Removed %q\n", removedWord)

	}
}

func addIgnoredFile(sc *SearchConfig) {
	// if ignored files in config is empty and there is a .gitignore in the current directory, ask if the user wants to add use that
	// Look at using filepath.Match or fs.Glob. Also research hot .gitnore files work
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
	check(err)

	providedIndices := strings.Split(indicesToRemove, " ")

	for i, indexString := range providedIndices {
		indexString = strings.TrimSpace(indexString)
		indexValue, err := strconv.ParseInt(indexString, 10, strconv.IntSize)
		check(err)

		indexValue -= 1 + int64(i) // because menu is 1-based and indices become progressively off by one more each time around the loop as an item is removed
		removedWord := sc.IgnoredFiles[indexValue]
		copy(sc.IgnoredFiles[indexValue:], sc.IgnoredFiles[indexValue+1:])
		sc.IgnoredFiles = sc.IgnoredFiles[:len(sc.IgnoredFiles)-1]

		fmt.Printf("Removed %q\n", removedWord)
	}
}

func getGitIgnorePatterns() []string {
	file, err := os.Open(".gitignore")
	utils.Check(err)
	defer utils.CloseFile(file)

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
