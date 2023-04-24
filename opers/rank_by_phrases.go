package opers

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"

	textrank "github.com/DavidBelicza/TextRank/v2"
	"github.com/DavidBelicza/TextRank/v2/rank"
)

func rankPhrases(inputFile string) ([]rank.SingleWord, error) {
	f, err := os.Open(inputFile)
	if err != nil {
		return nil, fmt.Errorf("Error os.Open: %s", err)
	}
	defer f.Close()

	readBytes, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("Error ioutil.ReadAll: %s", err)
	}

	tr := textrank.NewTextRank()
	tr.Populate(string(readBytes), textrank.NewDefaultLanguage(), textrank.NewDefaultRule())
	tr.Ranking(textrank.NewChainAlgorithm())
	phrases := textrank.FindPhrases(tr)

	var rankedPhrase []rank.SingleWord
	for _, phrase := range phrases {
		rankedPhrase = append(rankedPhrase, rank.SingleWord{Word: fmt.Sprintf("%s %s", phrase.Left, phrase.Right), Qty: phrase.Qty})
	}

	sort.Slice(rankedPhrase, func(i, j int) bool {
		return rankedPhrase[i].Qty > rankedPhrase[j].Qty
	})

	return rankedPhrase, nil
}
