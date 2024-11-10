package main

import (
	"context"
	"fmt"
	"os"

	"ceffo.com/bee/cmd"
)

func main() {
	ctx := context.Background()
	err := cmd.NewBeeCLI().Run(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
