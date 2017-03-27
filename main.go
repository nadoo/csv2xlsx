package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/tealeg/xlsx"
)

const (
	bom0 = 0xef
	bom1 = 0xbb
	bom2 = 0xbf
)

// BOMReader returns a Reader that discards the BOM header.
func BOMReader(r io.Reader) io.Reader {
	buf := bufio.NewReader(r)
	b, err := buf.Peek(3)
	if err != nil {
		return buf
	}
	if b[0] == bom0 && b[1] == bom1 && b[2] == bom2 {
		buf.Discard(3)
	}
	return buf
}

func csv2xlsx(csvPath string) {
	csvFile, err := os.Open(csvPath)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
	defer csvFile.Close()

	reader := csv.NewReader(BOMReader(csvFile))
	xlsxFile := xlsx.NewFile()
	sheetNo := 0
	sheetName := fmt.Sprintf("Sheet%d", sheetNo)
	sheet, _ := xlsxFile.AddSheet(sheetName)

	lineNum := 0
	fields, err := reader.Read()
	for err == nil {
		lineNum++
		// a sheet can contain 1048576 rows, 16384 columns.
		if lineNum%1000000 == 0 {
			sheetNo++
			sheetName = fmt.Sprintf("Sheet%d", sheetNo)
			sheet, _ = xlsxFile.AddSheet(sheetName)
		}

		row := sheet.AddRow()
		for _, field := range fields {
			cell := row.AddCell()
			cell.Value = field
		}

		fields, err = reader.Read()
	}
	if err != nil {
		fmt.Printf(err.Error())
	}

	fileName := strings.TrimSuffix(path.Base(csvPath), ".csv")
	outFile := fileName + ".xlsx"
	xlsxFile.Save(outFile)
}

func main() {
	if len(os.Args) > 1 {
		csv2xlsx(os.Args[1])
		return
	}

	csvFiles, _ := filepath.Glob("./*.csv")
	for _, csvFile := range csvFiles {
		csv2xlsx(csvFile)
	}
}
