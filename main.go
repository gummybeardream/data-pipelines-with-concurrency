package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {

	mode := flag.String(
		"mode",
		"concurrent",
		"processing mode: concurrent or sequential",
	)
	flag.Parse()

	imagePaths := []string{"images/corgi.jpeg",
		"images/bernesemountain.jpeg",
		"images/icelandicsheep.jpeg",
		"images/norwegianelkhound.jpeg",
	}

	results, err := processImages(*mode, imagePaths)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	hadError := false

	for _, result := range results {
		if result.Err != nil {
			hadError = true
			fmt.Fprintf(os.Stderr, "Failed: %v\n", result.Err)
			continue
		}

		fmt.Printf("Success: %s\n", result.OutPath)
	}
	if hadError {
		os.Exit(1)
	}
}
