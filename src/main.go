package main

import (
	"fmt"
)

func main() {
	fmt.Println("Start of loading.")
	var index, err = LoadIndex("./index.json")
	if err != nil {
		fmt.Println(err)
	}
	index.TUI()
}
