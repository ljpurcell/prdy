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

type Config struct {
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
	defaultConfigStruct := Config{[]string{"var_dump", "dd", "console.log"}, []string{"console.log(error)", "console.log(exception)"}}
	jsonConfigData, err := json.Marshal(defaultConfigStruct)
	check(err)
	os.WriteFile(".prdy_config.json", jsonConfigData, 0644)
}

func displayConfigMenu() {

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

	fmt.Println("Read config into memory")

	// Check whether the user is running or configuring the tool
	configFlagPtr := flag.Bool("config", false, "Bool (default: false). Open the configuration menu instead of immediately running the tool.")
	flag.Parse()

	if *configFlagPtr {
		fmt.Println("configuring")
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
	fmt.Printf("Value: %v. Type: %T\n", f, f)

	// if errors array isn't empty, display lines for each file and prompt user if they want to run tests or abort to fix

	// run tests if errors array is empty
}
