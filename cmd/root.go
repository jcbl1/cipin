package cmd

import (
	"context"
	"log"

	"github.com/jcbl1/cipin/opers"
	"github.com/spf13/cobra"
)

var (
	limitOpt              int
	inputFile, outputFile string
	debug                 bool
	noTrimmingOpt         bool
	phrasesOpt            bool
)

// Adding flags
func init() {
	rootCmd.Flags().IntVarP(&limitOpt, "limit", "l", 500, "Set limit of words in result")
	rootCmd.Flags().StringVarP(&inputFile, "input-file", "i", "", "Set the file to be processed")
	rootCmd.Flags().StringVarP(&outputFile, "output-file", "o", "output.xlsx", "File to store the result")
	rootCmd.Flags().BoolVarP(&debug, "debug", "d", false, "Debug mode")
	rootCmd.Flags().BoolVar(&noTrimmingOpt, "no-trimming", false, "No trimming of stop words and useless junk words")
	rootCmd.Flags().BoolVar(&phrasesOpt, "phrases", false, "Parse phrases instead of words")
}

// rootCmd will run when called without any subcommands
var rootCmd = &cobra.Command{
	Use:                   "cipin --input-file INPUT_FILE [-output-file OUTPUT_FILE] [--limit 200]",
	DisableFlagsInUseLine: true,
	Short:                 "cipin counts and ranks the words by frequency",
	Run: func(cmd *cobra.Command, args []string) {
		// Show usage if flag "input-file" is not specified
		if inputFile == "" {
			cmd.Usage()
			return
		}
		// Some initializations when in debug mode
		if debug {
			log.SetFlags(log.Llongfile)
		}

		// Process as users want
		err := process()
		if err != nil {
			log.Fatalln(err)
		}
	},
}

// Execute executes the commands
func Execute() {
	rootCmd.Execute()
}

// process calls opers.RankAndWrite and returns its output, if any
func process() error {
	// Create context in case the program doesn't exit as expected
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	return opers.RankAndWrite(ctx, inputFile, outputFile, noTrimmingOpt, limitOpt, phrasesOpt)
}
