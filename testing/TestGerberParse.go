package main

import (
	"os"
	"fmt"
	"gerber_rs274x"
	"path/filepath"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Error must give filename to parse as argument")
		os.Exit(1)
	}
	
	if inputFile,err := os.Open(os.Args[1]); err != nil {
		fmt.Printf("Error opening input file %s: %s\n", os.Args[1], err.Error())
		os.Exit(2)
	} else {
		
		if parsedFile,err := gerber_rs274x.ParseGerberFile(inputFile); err != nil {
			inputFile.Close()
			fmt.Printf("Error parsing gerber file: %v\n", err)
			os.Exit(3)
		} else {
			inputFile.Close()
			
			outputFileName := filepath.Base(os.Args[1] + ".png")
			
			if err := gerber_rs274x.GenerateSurface(outputFileName, parsedFile); err != nil {
				fmt.Printf("Error generating PNG file: %s\n", err.Error())
				os.Exit(5)
			}
		
			/*
			if outputFile,err := os.Create(filepath.Base(os.Args[1]) + ".svg"); err != nil {
				fmt.Printf("Error opening output file %s: %s\n", filepath.Base(os.Args[1]) + ".svg", err.Error())
				os.Exit(4) 
			} else {
				if err := gerber_rs274x.GenerateSurface(outputFile, parsedFile); err != nil {
					outputFile.Close()
					fmt.Printf("Error generating SVG file: %s\n", err.Error())
					os.Exit(5)
				} else {
					outputFile.Close()
				}
			}
			*/
		}
	}
}