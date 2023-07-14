package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"strconv"
	"strings"
)

/*
 * TODO: refactor adding and removing duplication in config functions
 */

// Options for optimising: pointers, replacements for concatenation, reuse variables
// google other options 'https://golangdocs.com/techniques-to-maximize-your-go-applications-performance'
// BENCHMARK first so that have something to write about

type SearchConfig struct {
	HitWords          []string
	ExcludedWords     []string
	SourceDirectories []string
	IgnoredFiles      []string
}

func (sc *SearchConfig) addToField(element string, field []string) {
	defer sc.updateConfigFile()
	field = append(field, element)
}

// Undecided whether this should take in string or slice of strings to make adding multiple fields easier
func (sc *SearchConfig) updateConfigFile() {
	scJson, err := json.Marshal(*sc)
	check(err)
	os.WriteFile(".prdy_config.json", scJson, 0644)
}

/*
 * PRIMARY FUNCTIONALITY
 */
func checkFileForHits(file fs.File, searchConfig *SearchConfig) {
	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	var outputMap = make(map[string][]string)
	var outputArray []string

	fileInfo, err := file.Stat()
	check(err)
	fileName := fileInfo.Name()

	for i := 1; fileScanner.Scan(); i++ {
		if slice_element_is_substring(searchConfig.HitWords, fileScanner.Text()) && !slice_element_is_substring(searchConfig.ExcludedWords, fileScanner.Text()) {
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

/*
 * CONFIGURATION
 */
func addNewHitWord(searchConfig *SearchConfig) {
	defer updateConfigFile(searchConfig)

	fmt.Println("\n\t* Adding new Hit Word *")
	fmt.Println("\nPlease type the word or pattern you want to match on.")
	fmt.Println("TIP: If you are looking for a function, leave off the parenthesis -- unless you know the exact naming of the argument(s) it has been called with.")
	fmt.Print("Add hit word: ")
	var newHitWord string
	fmt.Scanln(&newHitWord)

	searchConfig.HitWords = append(searchConfig.HitWords, newHitWord)

	fmt.Printf("Added %q\n", newHitWord)
	displayHitWords(searchConfig, false)
}

func removeHitWord(searchConfig *SearchConfig) {
	defer updateConfigFile(searchConfig)

	fmt.Println("\n\t* Remove Hit Word *")
	fmt.Println("\nPlease type the number of the word you want to remove.")
	fmt.Println("TIP: If you want to remove mutliple words, type a space seperated list.")

	displayHitWords(searchConfig, true)

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
		removedWord := searchConfig.HitWords[indexValue]
		copy(searchConfig.HitWords[indexValue:], searchConfig.HitWords[indexValue+1:])
		searchConfig.HitWords = searchConfig.HitWords[:len(searchConfig.HitWords)-1]

		fmt.Printf("Removed %q\n", removedWord)
	}
}

func addNewExcludedWord(searchConfig *SearchConfig) {
	defer updateConfigFile(searchConfig)

	fmt.Println("\n\t* Adding new Excluded Word *")
	fmt.Println("\nPlease type the word or pattern you want to PREVENT matching on.")
	fmt.Println("TIP: This is where you can use specific argument names to stop returning false positives.")
	fmt.Print("Add excluded word: ")
	var newExcludedWord string
	fmt.Scanln(&newExcludedWord)

	searchConfig.ExcludedWords = append(searchConfig.ExcludedWords, newExcludedWord)

	fmt.Printf("Added %q\n", newExcludedWord)
	displayExcludedVersions(searchConfig, false)
}

func removeExcludedWord(searchConfig *SearchConfig) {
	defer updateConfigFile(searchConfig)

	fmt.Println("\n\t* Remove Hit Word *")
	fmt.Println("\nPlease type the number of the word you want to remove.")
	fmt.Println("TIP: If you want to remove mutliple words, type a space seperated list.")

	displayExcludedVersions(searchConfig, true)

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
		removedWord := searchConfig.ExcludedWords[indexValue]
		copy(searchConfig.ExcludedWords[indexValue:], searchConfig.ExcludedWords[indexValue+1:])
		searchConfig.ExcludedWords = searchConfig.ExcludedWords[:len(searchConfig.ExcludedWords)-1]

		fmt.Printf("Removed %q\n", removedWord)

	}
}

func addIgnoredFile(searchConfig *SearchConfig) {
	// if ignored files in config is empty and there is a .gitignore in the current directory, ask if the user wants to add use that
	// Look at using filepath.Match or fs.Glob. Also research hot .gitnore files work
}

func removeIgnoredFile(searchConfig *SearchConfig) {

	defer updateConfigFile(searchConfig)

	fmt.Println("\n\t* Remove Ignored File *")
	fmt.Println("\nPlease type the number of the word you want to remove.")
	fmt.Println("TIP: If you want to remove multiple files, type a space seperated list.")

	displayIgnoredFiles(searchConfig, true)

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
		removedWord := searchConfig.IgnoredFiles[indexValue]
		copy(searchConfig.IgnoredFiles[indexValue:], searchConfig.IgnoredFiles[indexValue+1:])
		searchConfig.IgnoredFiles = searchConfig.IgnoredFiles[:len(searchConfig.IgnoredFiles)-1]

		fmt.Printf("Removed %q\n", removedWord)
	}
}

func runConfigWizard() {
	defaultConfigStruct := SearchConfig{[]string{"var_dump", "dd", "console.log"}, []string{"console.log(error)", "console.log(exception)"}, []string{}}
	jsonConfigData, err := json.Marshal(defaultConfigStruct)
	check(err)
	os.WriteFile(".prdy_config.json", jsonConfigData, 0644)
}

/*
 * MENU & USER INTERACTION
 */
func displayConfigMenu(searchConfig SearchConfig) {
	fmt.Println("\n\t--- CONFIG MENU ---")
	displayHitWords(&searchConfig, false)
	displayExcludedVersions(&searchConfig, false)
}

func displayHitWords(searchConfig *SearchConfig, displayIndices bool) {
	fmt.Println("\nYour Hit Words are:")
	if displayIndices {
		for i, v := range searchConfig.HitWords {
			fmt.Printf("\t%d. %s\n", (i + 1), v)
		}
	} else {
		for _, v := range searchConfig.HitWords {
			fmt.Printf("\t%s\n", v)
		}
	}
}

func displayExcludedVersions(searchConfig *SearchConfig, displayIndices bool) {
	fmt.Println("\nYour Excluded Versions are:")
	if displayIndices {
		for i, v := range searchConfig.ExcludedWords {
			fmt.Printf("\t%d. %s\n", i+1, v)
		}
	} else {
		for _, v := range searchConfig.ExcludedWords {
			fmt.Printf("\t%s\n", v)
		}
	}
}

func displayIgnoredFiles(searchConfig *SearchConfig, displayIndices bool) {

	fmt.Println("\nYour Ignored Files are:")
	if displayIndices {
		for i, v := range searchConfig.IgnoredFiles {
			fmt.Printf("\t%d. %s\n", i+1, v)
		}
	} else {
		for _, v := range searchConfig.IgnoredFiles {
			fmt.Printf("\t%s\n", v)
		}
	}
}

func getUserSelection() *int {
	// Refactor into sections for configuring hit words, then prompt for adding or removing
	menuOptions := []string{"\t1. Add hit word",
		"\t2. Remove hit word",
		"\t3. Add excluded version",
		"\t4. Remove excluded version",
		"\t5. Add ignored file or file type",
		"\t6. Remove ignored file or file type",
		"\t7. Stop configuring and run program",
		"\t8. Quit all together",
	}
	fmt.Println("\nWhat would you like to do?")
	for _, option := range menuOptions {
		fmt.Println(option)
	}
	fmt.Print("\nType the corresponding number and press enter: ")

	var selection int
	fmt.Scan(&selection)

	if selection > len(menuOptions) || selection == 0 {
		fmt.Println("\nSorry, I didn't get that...")
		getUserSelection()
	}

	return &selection
}

/*
 * FILE HANDLING
 */
func closeFile(file fs.File) {
	err := file.Close()
	check(err)
}

func readAndPrintFileByLine(file fs.File, searchConfig *SearchConfig) {
	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	var outputMap = make(map[string][]string)
	var outputArray []string

	fileInfo, err := file.Stat()
	check(err)
	fileName := fileInfo.Name()

	for i := 1; fileScanner.Scan(); i++ {
		if slice_element_is_substring(searchConfig.HitWords, fileScanner.Text()) && !slice_element_is_substring(searchConfig.ExcludedWords, fileScanner.Text()) {
			line := fmt.Sprintf("%v %v\n", i, fileScanner.Text())
			outputArray = append(outputArray, line)
		}
	}

	outputMap[fileName] = outputArray

	for _, line := range outputMap[fileName] {
		fmt.Printf(line)
	}

}

/*
 * ERROR HANDLING
 */
func check(err error) {
	if err != nil {
		panic(err)
	}
}

/*
 * UTILITIES
 */
func slice_element_is_substring(needles []string, haystack string) bool {
	for _, needle := range needles {
		if strings.Contains(haystack, needle) {
			return true
		}
	}
	return false
}

func getGitIgnorePatterns() []string {
	file, err := os.Open(".gitignore")
	check(err)
	defer closeFile(file)

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

func checkIfUserWantsToSetUpConfigFile() bool {
	fmt.Println("It appears you don't have a configuration file in this directory.")
	fmt.Println("You may be setting up a new project, or accidentally running this tool outside the base of your project.")
	fmt.Println("\nEnter y[es] if you'd like to create a new configuration file. (Entering n[o] will quit.)")

	var response string
	for {
		fmt.Scanln(&response)

		if strings.HasPrefix(response, "y") || strings.HasPrefix(response, "Y") {
			return true
		}

		if strings.HasPrefix(response, "n") || strings.HasPrefix(response, "N") {
			return false
		}
	}
}

func main() {
	runTool := true

	if _, err := os.Stat(".prdy_config.json"); err != nil {
		wantsToSetUpConfig := checkIfUserWantsToSetUpConfigFile()

		if wantsToSetUpConfig {
			runConfigWizard()
		}
	}

	configJson, err := os.ReadFile(".prdy_config.json")
	check(err)
	var searchConfig SearchConfig
	json.Unmarshal(configJson, &searchConfig)

	// Check whether the user is running or configuring the tool
	userIsConfiguring := flag.Bool("config", false, "Bool (default: false). Open the configuration menu instead of immediately running the tool.")
	flag.Parse()

	if *userIsConfiguring {

		runTool = false
		showMenu := true

		for showMenu {
			displayConfigMenu(searchConfig)
			userSelection := getUserSelection()

			switch *userSelection {
			case 1:
				addNewHitWord(&searchConfig)
			case 2:
				removeHitWord(&searchConfig)
			case 3:
				addNewExcludedWord(&searchConfig)
			case 4:
				removeExcludedWord(&searchConfig)
			case 5:
				addIgnoredFile(&searchConfig)
			case 6:
				removeIgnoredFile(&searchConfig)
			case 7:
				runTool = true
				showMenu = false
			case 8:
				showMenu = false
			}
		}
	}

	if runTool {
		// For getting the current working directory. Add basic checks and error handling if there isn't a .env file or app folder, ask the user if they are running it from the root
		pwd, err := os.Getwd()

		check(err)
		fsys := os.DirFS(pwd)
		patterns := getGitIgnorePatterns()

		// pattern := "*.go"
		// files, err := fs.Glob(fsys, pattern)
		// check(err)

		// Loop over provided source code directories, passing each to walkdDir
		// fs.WalkDir(fsys, ".", func(path string, directory fs.DirEntry, err error) error {
		//
		// 	for _, pattern := range patterns {
		// 		matched, err := filepath.Match(pattern, path)
		// 		check(err)
		//
		// 		if matched {
		// 			fmt.Printf("File %s matched for pattern %q\n", path, pattern)
		// 		}
		// 	}
		// 	return nil
		// })

		// if errors array isn't empty, display lines for each file and prompt user if they want to run tests or abort to fix

		// run tests if errors array is empty

	}
}
