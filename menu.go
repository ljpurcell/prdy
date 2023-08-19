package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"strings"
)

/*
 * MENU & USER INTERACTION
 */
func runConfigProcess(userHasConfigFile bool) {

	if userHasConfigFile {
		sc := loadConfig()
		displayConfigMenu(sc)
	} else {
		sc := createEmptyConfig()
		setUpConfigFile(sc)
	}
}

func setUpConfigFile(sc *SearchConfig) {
	fmt.Println("It appears you don't have a config file for this tool.")
	fmt.Println("Let's set one up now.")
	addHitWord(sc)
	addExcludedWord(sc)
	addIgnoredFile(sc)
}

func checkIfUserWantsToRunTool() bool {
	fmt.Println("\nGreat! You now have a config in place. Would you like to run the tool? Enter 'y' for yes, or anything else for no.")
	fmt.Print("Run tool: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	input := strings.ToLower(scanner.Text())

	if strings.HasPrefix(input, "y") {
		return true
	} else {
		return false
	}
}

func displayConfigMenu(sc *SearchConfig) {
	fmt.Println("\n\t--- CONFIG MENU ---")
	fmt.Println("What would you like to do?")

	menuOptions := []string{
		"1. Add to field",
		"2. Remove from field",
		"3. View your current configuration",
		"4. Quit and run tool",
		"5. Quit everything",
	}

	fmt.Print("\nType the corresponding number and press enter: ")
	for _, v := range menuOptions {
		fmt.Printf("\t %s\n", v)
	}

	menuSize := len(menuOptions)
	userSelection := getUserSelection(menuSize)

	switch *userSelection {
	case 1:
		displayAddToFieldOptions(sc)
	case 2:
		displayRemoveFromFieldOptions(sc)
	case 3:
		displayCurrentConfig(sc)
	case 4:
		runTool(sc)
	case 5:
		quitEverything()
	}
}

func getUserSelection(menuSize int) *int {
	var selection int
	fmt.Scan(&selection)

	for selection < 1 || selection > menuSize {
		fmt.Println("\nSorry, I didn't get that...")
		fmt.Scan(&selection)
	}

	return &selection
}

func displayAddToFieldOptions(sc *SearchConfig) {
	addMenuOptions := []string{
		"1. Add hit word",
		"2. Add excluded version",
		"3. Add ignored file",
	}
	menuSize := len(addMenuOptions)
	userSelection := getUserSelection(menuSize)

	switch *userSelection {
	case 1:
		addHitWord(sc)
	case 2:
		addExcludedWord(sc)
	case 3:
		addIgnoredFile(sc)
	}
}

func displayRemoveFromFieldOptions(sc *SearchConfig) {
	removeMenuOptions := []string{
		"1. Remove hit word",
		"2. Remove excluded version",
		"3. Remove ignored file",
	}
	menuSize := len(removeMenuOptions)
	userSelection := getUserSelection(menuSize)

	switch *userSelection {
	case 1:
		removeHitWord(sc)
	case 2:
		removeExcludedWord(sc)
	case 3:
		removeIgnoredFile(sc)
	}
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

func displayCurrentConfig(sc *SearchConfig) {
	displayHitWords(sc, false)
	displayExcludedVersions(sc, false)
	displayIgnoredFiles(sc, false)
}

func addHitWord(sc *SearchConfig) {
	fmt.Println("\n\t* Add Hit Words *")
	fmt.Println("\nPlease type the pattern you want to match on. To add multiple, use a space seperated list.")
	fmt.Println("TIP: If you are looking for a function, leave off the parenthesis -- unless you know the exact naming of the argument(s) it has been called with.")
	fmt.Print("Add hit word: ")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	err := scanner.Err()
	Check(err)

	wordsToAdd := strings.Split(scanner.Text(), " ")
	sc.addToField(wordsToAdd, &sc.HitWords)
}

func addExcludedWord(sc *SearchConfig) {
	fmt.Println("\n\t* Adding Excluded Words *")
	fmt.Println("\nPlease type the word or pattern you want to PREVENT matching on.")
	fmt.Println("TIP: This is where you can use specific argument names to stop returning false positives.")
	fmt.Print("Add excluded word: ")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	err := scanner.Err()
	Check(err)

	wordsToAdd := strings.Split(scanner.Text(), " ")
	sc.addToField(wordsToAdd, &sc.ExcludedWords)
}

func addIgnoredFile(sc *SearchConfig) {
	_, err := os.Stat(".prdy_config.json")
	userHasGitIgnoreFile := err == nil

	if userHasGitIgnoreFile {
		fmt.Println("\n\t* Add Ignored Files *")
		fmt.Println("\nYou have a .gitignore file in this directory. Would you like to automatically import it?")
		fmt.Println("This will prevent matching on anything within your ignored files. Use 'y' to confirm or nothing to skip.")
		fmt.Print("Import .gitignore file: ")

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		err = scanner.Err()
		Check(err)

		input := strings.ToLower(scanner.Text())

		if strings.HasPrefix(input, "y") {

			file, err := os.Open(".gitignore")
			defer file.Close()
			Check(err)

			scanner = bufio.NewScanner(file)
			scanner.Split(bufio.ScanLines)

			var ignoredFiles []string
			for scanner.Scan() {
				if !strings.HasPrefix(scanner.Text(), "#") {
					ignoredFiles = append(ignoredFiles, scanner.Text())
				}
			}

			sc.addToField(ignoredFiles, &sc.IgnoredFiles)
		}
	}
}

func removeHitWord(sc *SearchConfig) {
	fmt.Println("\n\t* Remove Hit Word *")
	fmt.Println("\nPlease type the number of the word you want to remove.")
	fmt.Println("TIP: If you want to remove mutliple words, type a space seperated list.")

	displayHitWords(sc, true)

	fmt.Print("\nRemove hit word: ")
	var indicesToRemove string
	inputReader := bufio.NewReader(os.Stdin)
	indicesToRemove, err := inputReader.ReadString('\n')
	Check(err)

	providedIndices := strings.Split(indicesToRemove, " ")

	sc.removeFromField(providedIndices, &sc.ExcludedWords)
}

func runTool(sc *SearchConfig) {
	fmt.Printf("Running with configuration: %v\n", sc)

	pwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting PWD in runTool: %v", err)
	}
	fsys := os.DirFS(pwd)

	// file, err := fsys.Open("prdy_test_file.txt")
	//if err != nil {
	// 	fmt.Printf("Error in run tool: %v", err)
	// }

	fs.WalkDir(fsys, ".", func(path string, directory fs.DirEntry, err error) error {
		file, err := os.Open(path)
		if err != nil {
			fmt.Printf("Error opening file in runTool: %v", path)
		}
		CheckFileForHits(file, sc)
		return nil
	})
	fmt.Println("Done!")
}

func quitEverything() {
	return
}
