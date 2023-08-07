package main

import(
    "bufio"
    "fmt"
    "strings"
    "os"
    "strconv"
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

func displayCurrentConfig(sc *SearchConfig){
    displayHitWords(sc, false)
    displayExcludedVersions(sc, false)
    displayIgnoredFiles(sc, false)
}


func addHitWord(sc *SearchConfig) {
	defer sc.updateConfigFile()

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


// Can be refactored, doesn't use sc.RemoveFromField method
func removeHitWord(sc *SearchConfig) {
    defer sc.updateConfigFile()

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

	for i, indexString := range providedIndices {
		indexString = strings.TrimSpace(indexString)
		indexValue, err := strconv.ParseInt(indexString, 10, strconv.IntSize)
		Check(err)

		indexValue -= 1 + int64(i) // because menu is 1-based and indices become progressively off by one more each time around the loop as an item is removed
		removedWord := sc.HitWords[indexValue]
		copy(sc.HitWords[indexValue:], sc.HitWords[indexValue+1:])
		sc.HitWords = sc.HitWords[:len(sc.HitWords)-1]

		fmt.Printf("Removed %q\n", removedWord)
	}
}

func runTool(sc *SearchConfig) {
    fmt.Printf("Running with configuration: %v\n", sc)
}

func quitEverything() {
    return
}
