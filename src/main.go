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
	fmt.Println("Start of saving.")
	err = index.Save()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Done")
}
