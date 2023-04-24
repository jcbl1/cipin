package opers

import (
	"fmt"
	"os"

	"github.com/DavidBelicza/TextRank/v2/rank"
	"github.com/xuri/excelize/v2"
)

func writeToExcel(rankedWords []rank.SingleWord, outputFile string, limitOpt int) error {
	f := excelize.NewFile()
	sheet := f.GetSheetName(0)

	for i, w := range rankedWords {
		if i >= limitOpt {
			break
		}
		f.SetCellStr(sheet, fmt.Sprintf("A%d", i+1), w.Word)
		f.SetCellInt(sheet, fmt.Sprintf("B%d", i+1), w.Qty)
	}

	outputF, err := os.OpenFile(outputFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("Error os.OpenFile: %s", err)
	}
	defer outputF.Close()

	_, err = f.WriteTo(outputF)
	if err != nil {
		return fmt.Errorf("Error f.WriteTo: %s", err)
	}

	return nil
}
