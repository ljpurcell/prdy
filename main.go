package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/rivo/tview"
)

var app = tview.NewApplication()
var pages = tview.NewPages()
var configForm = tview.NewForm()

func addConfigForm() {
	cfg := SearchConfig{}

	configForm.AddInputField("Hit Words", "", 20, nil, func(hw string) {
		hwSlice := strings.Split(hw, " ")
		cfg.HitWords = hwSlice
	})

	configForm.AddInputField("Excluded Words", "", 20, nil, func(ew string) {
		ewSlice := strings.Split(ew, " ")
		cfg.ExcludedWords = ewSlice
	})

	configForm.AddCheckbox("Import .gitignore file", false, func(importFile bool) {
		if importFile {
			getGitIgnorePatterns()
		}
	})
}

func main() {

	pages.AddPage("Config Form", configForm, true, true)

	/*
	 * Workflow:
	 * 1. checks branch is up to date - LAST
	 * 2. checks if user is configuring, configures if so
	 * 3. checks if there is a config file, configures if so
	 * 4. runs main tool for debugging matches
	 * 5. runs tests, output TBD
	 */

	// check if current branch is up to date with origin

	// Check whether the user is configuring the tool
	userIsConfiguring := flag.Bool("config", false, "Bool (default: false). Open the configuration menu instead of immediately running the tool.")
	flag.Parse()

	// If no config file, set up must be run
	_, err := os.Stat(".prdy_config.json")
	userHasConfigFile := err == nil
	wantsToRunTool := true

	if !userHasConfigFile || *userIsConfiguring {

		if err := app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
			panic(err)
		}

		pages.SwitchToPage("Config Form")

		runConfigProcess(userHasConfigFile)
		wantsToRunTool = checkIfUserWantsToRunTool()
	}

	if wantsToRunTool {
		sc := loadConfig()
		runTool(sc)
	}

	//	if runTool {
	// For getting the current working directory. Add basic utils.Checks and error handling if there isn't a .env file or app folder, ask the user if they are running it from the root

	// patterns := getGitIgnorePatterns()

	// pattern := "*.go"
	// files, err := fs.Glob(fsys, pattern)
	// utils.Check(err)

	// if errors array isn't empty, display lines for each file and prompt user if they want to run tests or abort to fix

	// run tests if errors array is empty

}

// }

// Options for optimising: pointers, replacements for concatenation, reuse variables
// google other options 'https://golangdocs.com/techniques-to-maximize-your-go-applications-performance'
// BENCHMARK first so that have something to write about

func ReadAndPrintFileByLine(file fs.File, sc *SearchConfig) {
	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	var outputMap = make(map[string][]string)
	var outputArray []string

	fileInfo, err := file.Stat()
	Check(err)
	fileName := fileInfo.Name()

	for i := 1; fileScanner.Scan(); i++ {
		if foundMatch(sc.HitWords, fileScanner.Text()) && !foundMatch(sc.ExcludedWords, fileScanner.Text()) {
			line := fmt.Sprintf("%v %v\n", i, fileScanner.Text())
			outputArray = append(outputArray, line)
		}
	}

	outputMap[fileName] = outputArray

	for _, line := range outputMap[fileName] {
		fmt.Printf(line)
	}

}
