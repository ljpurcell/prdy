package main

import (
	"fmt"
	"io/fs"
	"os"
)

func display(str, err string) {
	fmt.Println(str, err)
}

func main() {

	// pwd, err := os.Getwd()
	// fsys := os.DirFS(pwd)

	fsys := os.DirFS("/Users/LJPurcell/Code/prdy")

	f, err := fsys.Open("README.md")

	// fs.WalkDir(fsys,".", )

	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%T %v", f, f)
}
