package opers

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strings"

	textrank "github.com/DavidBelicza/TextRank/v2"
	"github.com/DavidBelicza/TextRank/v2/rank"
)

func rankWords(inputFile string) ([]rank.SingleWord, error) {
	f, err := os.Open(inputFile)
	if err != nil {
		return nil, fmt.Errorf("Error os.Open: %s", err)
	}
	defer f.Close()

	textBytes, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("Error ioutil.ReadAll: %s", err)
	}

	tr := textrank.NewTextRank()
	tr.Populate(string(textBytes), textrank.NewDefaultLanguage(), textrank.NewDefaultRule())
	tr.Ranking(textrank.NewChainAlgorithm())

	rankedWords := textrank.FindSingleWords(tr)
	sort.Slice(rankedWords, func(i, j int) bool {
		return rankedWords[i].Qty > rankedWords[j].Qty
	})

	return rankedWords, nil
}

func rankWordsNoTrimming(ctx context.Context, inputFile string) ([]rank.SingleWord, error) {
	f, err := os.Open(inputFile)
	if err != nil {
		return nil, fmt.Errorf("Error os.Open: %s", err)
	}
	defer f.Close()

	var word []byte
	wordsMap := make(map[string]int)
LOOP:
	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("Error ctx.Done() closed: %s", ctx.Err())
		default:
			var buf [1]byte
			_, err := f.Read(buf[:])
			if err == io.EOF {
				break LOOP
			}
			if err != nil {
				return nil, fmt.Errorf("Error f.Read: %s", err)
			}

			re := regexp.MustCompile(`\w{1}`)
			if re.Match(buf[:]) {
				word = append(word, buf[0])
			} else {
				if word != nil {
					wordsMap[strings.ToLower(string(word))]++
				}
				word = nil
			}
		}
	}

	var rankedWords []rank.SingleWord
	for k, v := range wordsMap {
		rankedWords = append(rankedWords, rank.SingleWord{Word: k, Qty: v})
	}

	return rankedWords, nil
}
