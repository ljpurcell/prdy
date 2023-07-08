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

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func closeFile(file fs.File) {
	err := file.Close()
	fmt.Println("Closed file")
	check(err)
}

func createConfigFile() {
	defaultConfigStruct := SearchConfig{[]string{"var_dump", "dd", "console.log"}, []string{"console.log(error)", "console.log(exception)"}}
	jsonConfigData, err := json.Marshal(defaultConfigStruct)
	check(err)
	os.WriteFile(".prdy_config.json", jsonConfigData, 0644)
}

func displayConfigMenu(searchConfig SearchConfig) {
	fmt.Println("\n\t--- CONFIG MENU ---")

	fmt.Println("\nHit words:")
	for _, v := range searchConfig.HitWords {
		fmt.Println("\t" + v)
	}

	fmt.Println("\nExcluded versions:")
	for _, v := range searchConfig.ExcludedVersion {
		fmt.Println("\t" + v)
	}
}

func getUserSelection() {
	menuOptions := []string{"\t1. Add hit word",
		"\t2. Remove hit word",
		"\t3. Add excluded version",
		"\t4. Remove excluded version",
		"\t5. Quit",
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

	switch selection {
	case 1:
		fmt.Println("user selected 1")
	case 2:
		fmt.Println("user selected 2")
	}
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
		getUserSelection()
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
