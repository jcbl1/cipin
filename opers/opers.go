package opers

import (
	"context"

	"github.com/DavidBelicza/TextRank/v2/rank"
)

// RankAndWrite takes parameters and use them to do the job as the function name tells
func RankAndWrite(ctx context.Context, inputFile string, outputFile string, noTrimmingOpt bool, limitOpt int, phrasesOpt bool) error {
	// Initialize "words" to be written and error to be returned
	var rankedWords []rank.SingleWord
	var err error

	// Assign rankedWords based on different cases
	if phrasesOpt {
		rankedWords, err = rankPhrases(inputFile)
	} else if noTrimmingOpt {
		rankedWords, err = rankWordsNoTrimming(ctx, inputFile)
	} else {
		rankedWords, err = rankWords(inputFile)
	}
	if err != nil {
		return err
	}

	// Call a function to write results to file and return its error, if any
	return writeToExcel(rankedWords, outputFile, limitOpt)
}
