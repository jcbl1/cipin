package opers

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"

	textrank "github.com/DavidBelicza/TextRank/v2"
	"github.com/DavidBelicza/TextRank/v2/rank"
)

// rankPhrases takes data in inputFile and returns a ranked phrases slice named "rankedWords"
func rankPhrases(inputFile string) ([]rank.SingleWord, error) {
	// Open file
	f, err := os.Open(inputFile)
	if err != nil {
		return nil, fmt.Errorf("Error os.Open: %s", err)
	}
	defer f.Close()

	// read all from f as bytes
	readBytes, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("Error ioutil.ReadAll: %s", err)
	}

	// processing the ranking
	tr := textrank.NewTextRank()
	tr.Populate(string(readBytes), textrank.NewDefaultLanguage(), textrank.NewDefaultRule())
	tr.Ranking(textrank.NewChainAlgorithm())
	// Getting ranked phrases (by weight)
	phrases := textrank.FindPhrases(tr)

	// create a SingleWord slice to fit the return's need
	var rankedPhrase []rank.SingleWord
	// copying phrases to rankedPhrase
	for _, phrase := range phrases {
		rankedPhrase = append(rankedPhrase, rank.SingleWord{Word: fmt.Sprintf("%s %s", phrase.Left, phrase.Right), Qty: phrase.Qty})
	}

	// Sorting decsending order
	sort.Slice(rankedPhrase, func(i, j int) bool {
		return rankedPhrase[i].Qty > rankedPhrase[j].Qty
	})

	return rankedPhrase, nil
}
