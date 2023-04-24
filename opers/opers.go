package opers

import (
	"context"

	"github.com/DavidBelicza/TextRank/v2/rank"
)

func RankAndWrite(ctx context.Context, inputFile string, outputFile string, noTrimmingOpt bool, limitOpt int, phrasesOpt bool) error {
	var rankedWords []rank.SingleWord
	var err error
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

	return writeToExcel(rankedWords, outputFile, limitOpt)
}
