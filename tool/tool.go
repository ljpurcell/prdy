package tool

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
