package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"strings"
)

// Options for optimising: pointers, printf over concatenation?, google other options 'https://golangdocs.com/techniques-to-maximize-your-go-applications-performance'
// BENCHMARK first so that have something to write about

type SearchConfig struct {
	HitWords        []string
	ExcludedVersion []string
}

/*
 * CONFIGURATION
 */
func actionUserSelection(userSelection int, searchConfig SearchConfig) {
	defer updateConfigFile(&searchConfig)

	switch userSelection {
	case 1:
		addNewHitWord(&searchConfig)
	case 2:
		removeHitWord(&searchConfig)
	}

}

func addNewHitWord(searchConfig *SearchConfig) {
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
	fmt.Println("\n\t* Remove Hit Word *")
	fmt.Println("\nPlease type the number of the word you want to remove.")
	fmt.Println("TIP: If you want to remove mutliple words, type a space seperated list.")

	displayHitWords(searchConfig, true)

	fmt.Print("\nRemove hit word: ")
	var indexToRemove int
	fmt.Scanln(&indexToRemove)
	indexToRemove-- // Menu is 1-based, so need to decrement index

	removedWord := searchConfig.HitWords[indexToRemove]
	copy(searchConfig.HitWords[indexToRemove:], searchConfig.HitWords[indexToRemove+1:])
	searchConfig.HitWords = searchConfig.HitWords[:len(searchConfig.HitWords)-1]

	fmt.Printf("Remove %q\n", removedWord)
}

func createConfigFile() {
	defaultConfigStruct := SearchConfig{[]string{"var_dump", "dd", "console.log"}, []string{"console.log(error)", "console.log(exception)"}}
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
	displayExcludedVersions(&searchConfig)
}

func displayHitWords(searchConfig *SearchConfig, displayIndices bool) {
	fmt.Println("\nYour Hit Words are:")

	if displayIndices {
		for i, v := range searchConfig.HitWords {
			fmt.Printf("\t%d. %s\n", i+1, v)
		}

	} else {
		for _, v := range searchConfig.HitWords {
			fmt.Println("\t" + v)
		}

	}
}

func displayExcludedVersions(searchConfig *SearchConfig) {
	fmt.Println("\nYour Excluded Versions are:")
	for _, v := range searchConfig.ExcludedVersion {
		fmt.Println("\t" + v)
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
	fmt.Println("Closed file")
	check(err)
}

func readAndPrintFileByLine(file fs.File) {
	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	var outputMap = make(map[string][]string)
	var outputArray []string

	fileInfo, err := file.Stat()
	check(err)
	fileName := fileInfo.Name()

	for i := 1; fileScanner.Scan(); i++ {
		if strings.Contains(fileScanner.Text(), "#") {
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

func main() {

	if _, err := os.Stat(".prdy_config.json"); err != nil {
		createConfigFile()
	}

	configJson, err := os.ReadFile(".prdy_config.json")
	var searchConfig SearchConfig
	json.Unmarshal(configJson, &searchConfig)

	// Check whether the user is running or configuring the tool
	userIsConfiguring := flag.Bool("config", false, "Bool (default: false). Open the configuration menu instead of immediately running the tool.")
	flag.Parse()

	if *userIsConfiguring {
		displayConfigMenu(searchConfig)
		userSelection := getUserSelection()
		actionUserSelection(*userSelection, searchConfig)
	} else {
		fmt.Println("Not configuring")
	}
	// if configuring, display current config and give option to add or remove terms from the 'hit' or 'exclude' lists
	// Will probably need to create a .prdy_config file in the directory to preseve config overtime

	// For getting the current working directory. Add basic checks and error handling if there isn't a .env file or app folder, ask the user if they are running it from the root
	// pwd, err := os.Getwd()
	// fsys := os.DirFS(pwd)

	// For walking the file system from root. Replace anonymous function with one that actually implements desired functionality
	// fs.WalkDir(fsys, ".", func(path string, directory fs.DirEntry, err error) error {
	// 	fmt.Println(path)
	// 	return nil
	// })

	fsys := os.DirFS("/Users/LJPurcell/Code/prdy")

	f, err := fsys.Open("README.md")
	check(err)
	defer closeFile(f)

	readAndPrintFileByLine(f)

	// if errors array isn't empty, display lines for each file and prompt user if they want to run tests or abort to fix

	// run tests if errors array is empty
}
