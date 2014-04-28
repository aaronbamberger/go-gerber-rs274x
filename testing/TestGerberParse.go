package main

import (
	"os"
	"fmt"
	"gerber_rs274x"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Error must give filename to parse as argument")
		os.Exit(1)
	}
	
	if file,err := os.Open(os.Args[1]); err != nil {
		fmt.Printf("Error opening given file %s: %v\n", os.Args[1], err)
		os.Exit(2)
	} else {
		
		if _,err := gerber_rs274x.ParseGerberFile(file); err != nil {
			file.Close()
			fmt.Printf("Error parsing gerber file: %v\n", err)
			os.Exit(3)
		}
		
		file.Close()
	}
	
	/*
	//input := "";
	
	//regex := regexp.MustCompile("X?(-?[[:digit:]]*)Y?(-?[[:digit:]]*)I?(-?[[:digit:]]*)J?(-?[[:digit:]]*)")
	
	
	
	results := regex.FindAllStringSubmatch(input, -1)
	
	fmt.Printf("Results Length: %d\n", len(results));
	for index, submatch := range results {
		fmt.Printf("Results(%d) Length: %d\n", index, len(submatch));
	}
	fmt.Printf("Results: %v\n", results)
	*/
}