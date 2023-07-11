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
	HitWords         []string
	ExcludedVersions []string
	IgnoredFiles     []string
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
		if slice_contains_substring(searchConfig.HitWords, fileScanner.Text()) && !slice_contains_substring(searchConfig.ExcludedVersions, fileScanner.Text()) {
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

	searchConfig.ExcludedVersions = append(searchConfig.ExcludedVersions, newExcludedWord)

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
		removedWord := searchConfig.ExcludedVersions[indexValue]
		copy(searchConfig.ExcludedVersions[indexValue:], searchConfig.ExcludedVersions[indexValue+1:])
		searchConfig.ExcludedVersions = searchConfig.ExcludedVersions[:len(searchConfig.ExcludedVersions)-1]

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

func createConfigFile() {
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
		for i, v := range searchConfig.ExcludedVersions {
			fmt.Printf("\t%d. %s\n", i+1, v)
		}
	} else {
		for _, v := range searchConfig.ExcludedVersions {
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
	menuOptions := []string{"\t1. Add hit word",
		"\t2. Remove hit word",
		"\t3. Add excluded version",
		"\t4. Remove excluded version",
		"\t5. Stop configuring and run program",
		"\t6. Quit all together",
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

func updateConfigFile(searchConfig *SearchConfig) {
	searchConfigJson, err := json.Marshal(*searchConfig)
	check(err)
	os.WriteFile(".prdy_config.json", searchConfigJson, 0644)
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
		if slice_contains_substring(searchConfig.HitWords, fileScanner.Text()) && !slice_contains_substring(searchConfig.ExcludedVersions, fileScanner.Text()) {
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
func slice_contains_substring(needles []string, haystack string) bool {
	for _, needle := range needles {
		if strings.Contains(haystack, needle) {
			return true
		}
	}
	return false
}

func main() {
	runTool := true

	if _, err := os.Stat(".prdy_config.json"); err != nil {
		createConfigFile()
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
				runTool = true
				showMenu = false
			case 6:
				showMenu = false
			}
		}
	}

	if runTool {
		// For getting the current working directory. Add basic checks and error handling if there isn't a .env file or app folder, ask the user if they are running it from the root
		// pwd, err := os.Getwd()
		// fsys := os.DirFS(pwd)

		fsys := os.DirFS("/Users/LJPurcell/Code/tragics/")

		// For walking the file system from root. Replace anonymous function with one that actually implements desired functionality
		fs.WalkDir(fsys, ".", func(path string, directory fs.DirEntry, err error) error {
			f, err := fsys.Open(path)
			check(err)
			defer closeFile(f)
			checkFileForHits(f, &searchConfig)
			return nil
		})

		// if errors array isn't empty, display lines for each file and prompt user if they want to run tests or abort to fix

		// run tests if errors array is empty

	}
}
