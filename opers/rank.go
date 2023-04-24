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

// rankWords reads data from inputFile and returns the words in it ranked by occurences of each one
func rankWords(inputFile string) ([]rank.SingleWord, error) {
	// Open file
	f, err := os.Open(inputFile)
	if err != nil {
		return nil, fmt.Errorf("Error os.Open: %s", err)
	}
	defer f.Close()

	// read all bytes into textBytes from f
	textBytes, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("Error ioutil.ReadAll: %s", err)
	}

	// Processing the data
	tr := textrank.NewTextRank()
	tr.Populate(string(textBytes), textrank.NewDefaultLanguage(), textrank.NewDefaultRule())
	tr.Ranking(textrank.NewChainAlgorithm()) // Use textrank.NewChainAlgorithm() because it is based on occurences

	// Getting ranked words
	rankedWords := textrank.FindSingleWords(tr)
	sort.Slice(rankedWords, func(i, j int) bool {
		return rankedWords[i].Qty > rankedWords[j].Qty
	})

	return rankedWords, nil
}

// rankWordsNoTrimming does the same job as rankWords except that it won't get rid of stop words and junk words
func rankWordsNoTrimming(ctx context.Context, inputFile string) ([]rank.SingleWord, error) {
	// Open file
	f, err := os.Open(inputFile)
	if err != nil {
		return nil, fmt.Errorf("Error os.Open: %s", err)
	}
	defer f.Close()

	// Getting a map, of which the key is the word and the value is the occurence of it
	var word []byte                  // word is the buffer to store a whole word
	wordsMap := make(map[string]int) //wordsMap is used to count occurences of each word

LOOP:
	for {
		select {
		case <-ctx.Done(): // In case this loop cannot stop correctly
			return nil, fmt.Errorf("Error ctx.Done() closed: %s", ctx.Err())
		default:
			// Read 1 byte
			var buf [1]byte
			_, err := f.Read(buf[:])
			// if true break the outter loop
			if err == io.EOF {
				break LOOP
			}
			if err != nil {
				return nil, fmt.Errorf("Error f.Read: %s", err)
			}

			// determine whether that single byte is an alphabetical one
			re := regexp.MustCompile(`\w{1}`)
			// if true, append it to the buffered word. If not, the wordsMap with the key of word adds 1
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

	// copying words from wordsMap to rankedWords and sort them
	var rankedWords []rank.SingleWord
	for k, v := range wordsMap {
		rankedWords = append(rankedWords, rank.SingleWord{Word: k, Qty: v})
	}

	// TODO: rankedWords aren't sorted yet!!!

	return rankedWords, nil
}
