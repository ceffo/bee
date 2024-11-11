package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"ceffo.com/bee/bee"
	"ceffo.com/bee/core"
)

const (
	FlagLetters = "letters"
)

func newSolveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "solve",
		Short: "Solve a bee spelling game",
		RunE:  solveCmd,
	}

	cmd.Flags().StringP(FlagLetters, "l", "", "Letters to use, starting with the center letter (required)")

	return cmd
}

func solveCmd(cmd *cobra.Command, _ []string) error {
	letters, err := cmd.Flags().GetString(FlagLetters)
	if err != nil {
		return err
	}
	// make an input out of the letters
	input, err := bee.NewFrom(letters)
	if err != nil {
		return err
	}
	worlistFile, err := cmd.InheritedFlags().GetString(FlagWordlist)
	if err != nil {
		return err
	}

	c, err := core.New(
		core.WithSourceMaker(worlistFile),
		core.WithStdoutLogging(),
	)
	if err != nil {
		return err
	}
	defer c.Close()
	solver := bee.NewSolver(c.Source())
	res := solver.SolveFor(input)
	words := make([]string, 0)
	for word := range res {
		words = append(words, word)
	}
	const wordsPerLine = 15
	for i := 0; i < len(words); i += wordsPerLine {
		end := i + wordsPerLine
		if end > len(words) {
			end = len(words)
		}
		fmt.Println(strings.Join(words[i:end], " "))
	}

	return nil
}
