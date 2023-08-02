package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"strings"
    "github.com/ljpurcell/prdy/utils"
    "github.com/ljpurcell/prdy/config"
)

func main() {
	runTool := true

	if _, err := os.Stat(".prdy_config.json"); err != nil {
		wantsToSetUpConfig := CheckIfUserWantsToSetUpConfigFile()

		if wantsToSetUpConfig {
			runConfigWizard()
		}
	}

	configJson, err := os.ReadFile(".prdy_config.json")
	utils.Check(err)
	var sc config.SearchConfig
	json.Unmarshal(configJson, &sc)

	// utils.Check whether the user is running or configuring the tool
	userIsConfiguring := flag.Bool("config", false, "Bool (default: false). Open the configuration menu instead of immediately running the tool.")
	flag.Parse()

	if *userIsConfiguring {

		runTool = false
		showMenu := true

		for showMenu {
			displayConfigMenu(sc)
			userSelection := getUserSelection()

			switch *userSelection {
			case 1:
				config.AddHitWords(&sc)
			case 2:
				config.RemoveHitWord(&sc)
			case 3:
				confing.AddExcludedWords(&sc)
			case 4:
				config.RemoveExcludedWord(&sc)
			case 5:
				config,AddIgnoredFile(&sc)
			case 6:
				config.RemoveIgnoredFile(&sc)
			case 7:
				runTool = true
				showMenu = false
			case 8:
				showMenu = false
			}
		}
	}

	if runTool {
		// For getting the current working directory. Add basic utils.Checks and error handling if there isn't a .env file or app folder, ask the user if they are running it from the root
		// pwd, err := os.Getwd()

		utils.Check(err)
		// fsys := os.DirFS(pwd)
		// patterns := getGitIgnorePatterns()

		// pattern := "*.go"
		// files, err := fs.Glob(fsys, pattern)
		// utils.Check(err)

		// Loop over provided source code directories, passing each to walkdDir
		// fs.WalkDir(fsys, ".", func(path string, directory fs.DirEntry, err error) error {
		//
		// 	for _, pattern := range patterns {
		// 		matched, err := filepath.Match(pattern, path)
		// 		utils.Check(err)
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

// Options for optimising: pointers, replacements for concatenation, reuse variables
// google other options 'https://golangdocs.com/techniques-to-maximize-your-go-applications-performance'
// BENCHMARK first so that have something to write about


/*
 * PRIMARY FUNCTIONALITY
 */
func CheckFileForHits(file fs.File, sc *SearchConfig) {
	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	var outputMap = make(map[string][]string)
	var outputArray []string

	fileInfo, err := file.Stat()
	utils.Check(err)
	fileName := fileInfo.Name()

	for i := 1; fileScanner.Scan(); i++ {
		if utils.SliceElementIsSubstring(sc.HitWords, fileScanner.Text()) && utils.SliceElementIsSubstring(sc.ExcludedWords, fileScanner.Text()) {
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


func runConfigWizard(sc *SearchConfig) {
	addHitWords(sc)
	addExcludedWords(sc)
	// prompt read in latest gitignore file
	os.WriteFile(".prdy_config.json", jsonConfigData, 0644)
}

/*
 * MENU & USER INTERACTION
 */
func displayConfigMenu(sc SearchConfig) {
	fmt.Println("\n\t--- CONFIG MENU ---")
	displayHitWords(&sc, false)
	displayExcludedVersions(&sc, false)
}

func displayHitWords(sc *SearchConfig, displayIndices bool) {
	fmt.Println("\nYour Hit Words are:")
	if displayIndices {
		for i, v := range sc.HitWords {
			fmt.Printf("\t%d. %s\n", (i + 1), v)
		}
	} else {
		for _, v := range sc.HitWords {
			fmt.Printf("\t%s\n", v)
		}
	}
}

func displayExcludedVersions(sc *SearchConfig, displayIndices bool) {
	fmt.Println("\nYour Excluded Versions are:")
	if displayIndices {
		for i, v := range sc.ExcludedWords {
			fmt.Printf("\t%d. %s\n", i+1, v)
		}
	} else {
		for _, v := range sc.ExcludedWords {
			fmt.Printf("\t%s\n", v)
		}
	}
}

func displayIgnoredFiles(sc *SearchConfig, displayIndices bool) {

	fmt.Println("\nYour Ignored Files are:")
	if displayIndices {
		for i, v := range sc.IgnoredFiles {
			fmt.Printf("\t%d. %s\n", i+1, v)
		}
	} else {
		for _, v := range sc.IgnoredFiles {
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




func CheckIfUserWantsToSetUpConfigFile() bool {
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

func ReadAndPrintFileByLine(file fs.File, sc *config.SearchConfig) {
	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	var outputMap = make(map[string][]string)
	var outputArray []string

	fileInfo, err := file.Stat()
	Check(err)
	fileName := fileInfo.Name()

	for i := 1; fileScanner.Scan(); i++ {
		if SliceElementIsSubstring(sc.HitWords, fileScanner.Text()) && !SliceElementIsSubstring(sc.ExcludedWords, fileScanner.Text()) {
			line := fmt.Sprintf("%v %v\n", i, fileScanner.Text())
			outputArray = append(outputArray, line)
		}
	}

	outputMap[fileName] = outputArray

	for _, line := range outputMap[fileName] {
		fmt.Printf(line)
	}

}
