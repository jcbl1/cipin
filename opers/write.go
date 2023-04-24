package opers

import (
	"fmt"
	"os"

	"github.com/DavidBelicza/TextRank/v2/rank"
	"github.com/xuri/excelize/v2"
)

// writeToExcel writes rankedWords to outputFile. It only writes up to limitOpt rows
func writeToExcel(rankedWords []rank.SingleWord, outputFile string, limitOpt int) error {
	// Creating new excelize.File
	f := excelize.NewFile()
	// Getting the first sheet name
	sheet := f.GetSheetName(0)

	// Write data to f
	for i, w := range rankedWords {
		if i >= limitOpt {
			break
		}
		f.SetCellStr(sheet, fmt.Sprintf("A%d", i+1), w.Word)
		f.SetCellInt(sheet, fmt.Sprintf("B%d", i+1), w.Qty)
	}

	// Open outputFile in readonly mode and Truncates it. If not existing, will create one
	outputF, err := os.OpenFile(outputFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("Error os.OpenFile: %s", err)
	}
	defer outputF.Close()

	// write buffered excelize.File to outputFile
	_, err = f.WriteTo(outputF)
	if err != nil {
		return fmt.Errorf("Error f.WriteTo: %s", err)
	}

	return nil
}
