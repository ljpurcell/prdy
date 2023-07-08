package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
)

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

func readAndPrintFileByLine(file fs.File) {
	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	for i := 1; fileScanner.Scan(); i++ {
		fmt.Printf("%v %v\n", i, fileScanner.Text())
	}

}
func main() {

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
}
