package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/xuri/excelize/v2"
)

var (
	limitOpt              int
	inputFile, outputFile string
	debug                 bool
)

func init() {
	rootCmd.Flags().IntVarP(&limitOpt, "limit", "l", 500, "Set limit of words in result")
	rootCmd.Flags().StringVarP(&inputFile, "input-file", "i", "", "Set the file to be processed")
	rootCmd.Flags().StringVarP(&outputFile, "output-file", "o", "output.xlsx", "File to store the result")
	rootCmd.Flags().BoolVarP(&debug, "debug", "d", false, "Debug mode")
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
	words, err := getAllWords()
	if err != nil {
		return err
	}

	originMap := cipin(words)

	return writeToExcel(originMap)
}

func getAllWords() ([]string, error) {
	f, err := os.Open(inputFile)
	if err != nil {
		return nil, fmt.Errorf("Error os.Open: %s", err)
	}
	defer f.Close()

	var word []byte
	var words []string
	for {
		var buf [1]byte
		_, err := f.Read(buf[:])
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("Error f.Read: %s", err)
		}

		re := regexp.MustCompile(`[a-zA-Z]{1}`)
		if re.Match(buf[:]) {
			word = append(word, buf[0])
		} else {
			if word != nil {
				words = append(words, string(word))
			}
			word = nil
		}
	}

	return words, nil
}

func cipin(words []string) map[string]int {
	ret := make(map[string]int)
	for _, word := range words {
		ret[strings.ToLower(word)]++
	}
	return ret
}

type CipinMap []struct {
	Word string
	Freq int
}

func writeToExcel(originMap map[string]int) error {
	var cipinMap CipinMap
	for k, v := range originMap {
		cipinMap = append(cipinMap, struct {
			Word string
			Freq int
		}{k, v})
	}

	sort.Slice(cipinMap, func(i, j int) bool {
		return cipinMap[i].Freq < cipinMap[j].Freq
	})

	f := excelize.NewFile()
	sheet := f.GetSheetName(0)

	n := len(cipinMap)
	for i := n - 1; i >= 0; i-- {
		j := n - 1 - i
		if j > limitOpt {
			break
		}
		f.SetCellStr(sheet, fmt.Sprintf("A%d", j+1), cipinMap[i].Word)
		f.SetCellInt(sheet, fmt.Sprintf("B%d", j+1), cipinMap[i].Freq)
	}

	outputFile, err := os.OpenFile(outputFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("Error os.OpenFile: %s", err)
	}
	_, err = f.WriteTo(outputFile)
	if err != nil {
		log.Println("here")
		return fmt.Errorf("Error f.WriteTo: %s", err)
	}

	return nil
}
