package cmd

import (
	"context"

	"github.com/spf13/cobra"

	"ceffo.com/bee/app"
	"ceffo.com/bee/core"
)

const (
	FlagWordlist = "wordlist"
	FlagLogfile  = "log"
)

type BeeCLI struct {
	rootCmd *cobra.Command
}

func NewBeeCLI() *BeeCLI {
	r := newRootCmd()
	r.AddCommand(
		newSolveCmd(),
	)
	return &BeeCLI{r}
}

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "bee",
		Short:         "Bee is a solver for the Spelling Bee puzzle",
		RunE:          handleRootCmd,
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.PersistentFlags().StringP(FlagWordlist, "w", "data/en.txt", "path to the word list file")
	cmd.PersistentFlags().String(FlagLogfile, "bee.log", "path to the log file")

	return cmd
}

func (b *BeeCLI) Run(ctx context.Context) error {
	return b.rootCmd.ExecuteContext(ctx)
}

func handleRootCmd(cmd *cobra.Command, _ []string) error {
	// no command, launch the TUI
	wordlistFile, err := cmd.Flags().GetString(FlagWordlist)
	if err != nil {
		return err
	}
	logFile, err := cmd.Flags().GetString(FlagLogfile)
	if err != nil {
		return err
	}

	c, err := core.New(
		core.WithFileLogging(logFile),
		core.WithSourceMaker(wordlistFile),
	)
	if err != nil {
		return err
	}
	defer c.Close()
	a := app.New(c)
	return a.Run()
}
