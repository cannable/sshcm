package main

import "fmt"

func add(cnArgs map[string]*string) {
	for f, s := range cnArgs {
		fmt.Println(f, ":", *s)
	}

	// TODO: Check that connection nickname starts with a letter

	//addCmd.PrintDefaults()

}

func set(id int, cnArgs map[string]*string) {
	fmt.Println("Set connection.")
	for f, s := range cnArgs {
		fmt.Println(f, ":", *s)
	}
	fmt.Println("id", ":", id)

	// TODO: Check that connection nickname starts with a letter

	//addCmd.PrintDefaults()

}
