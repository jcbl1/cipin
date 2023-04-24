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

func init() {
	rootCmd.Flags().IntVarP(&limitOpt, "limit", "l", 500, "Set limit of words in result")
	rootCmd.Flags().StringVarP(&inputFile, "input-file", "i", "", "Set the file to be processed")
	rootCmd.Flags().StringVarP(&outputFile, "output-file", "o", "output.xlsx", "File to store the result")
	rootCmd.Flags().BoolVarP(&debug, "debug", "d", false, "Debug mode")
	rootCmd.Flags().BoolVar(&noTrimmingOpt, "no-trimming", false, "No trimming of stop words and useless junk words")
	rootCmd.Flags().BoolVar(&phrasesOpt, "phrases", false, "Parse phrases instead of words")
}

var rootCmd = &cobra.Command{
	Use:                   "cipin --input-file INPUT_FILE [-output-file OUTPUT_FILE] [--limit 200]",
	DisableFlagsInUseLine: true,
	Short:                 "cipin counts and ranks the words by frequency",
	Run: func(cmd *cobra.Command, args []string) {
		if inputFile == "" {
			cmd.Usage()
			return
		}
		if debug {
			log.SetFlags(log.Llongfile)
		}

		err := process()
		if err != nil {
			log.Fatalln(err)
		}
	},
}

func Execute() {
	rootCmd.Execute()
}

func process() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	return opers.RankAndWrite(ctx, inputFile, outputFile, noTrimmingOpt, limitOpt, phrasesOpt)
}
