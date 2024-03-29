package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/tealeg/xlsx/v3"
)

const (
	bom0 = 0xef
	bom1 = 0xbb
	bom2 = 0xbf
)

var version = "1.3.0"

// BOMReader returns a Reader that discards the BOM header.
func BOMReader(ir io.Reader) io.Reader {
	r := bufio.NewReader(ir)
	b, err := r.Peek(3)
	if err != nil {
		return r
	}
	if b[0] == bom0 && b[1] == bom1 && b[2] == bom2 {
		r.Discard(3)
	}
	return r
}

func csv2xlsx(csvPath string) {
	fmt.Printf("\nProcessing: %s", csvPath)

	csvFile, err := os.Open(csvPath)
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	defer csvFile.Close()

	reader := csv.NewReader(BOMReader(csvFile))

	xlsxFile := xlsx.NewFile()
	xlsx.SetDefaultFont(10, "Verdana")

	sheetNo := 0
	sheetName := fmt.Sprintf("Sheet%d", sheetNo)
	sheet, _ := xlsxFile.AddSheet(sheetName)

	lineNo := 0
	fields, err := reader.Read()
	for err == nil {
		lineNo++
		// a sheet can contain 1048576 rows, 16384 columns.
		if lineNo%1000000 == 0 {
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

	if err != nil && err != io.EOF {
		fmt.Print(err.Error())
		return
	}

	fileName := strings.TrimSuffix(filepath.Base(csvPath), ".csv")
	outFile := fileName + ".xlsx"
	xlsxFile.Save(outFile)

	fmt.Printf(" -> %s", outFile)
}

func showHelp() {
	fmt.Println()
	fmt.Println("-")
	fmt.Println("Usage: csv2xlsx")
	fmt.Println("       csv2xlsx [FILEPATH]")
	fmt.Println()
	fmt.Println("Example:")
	fmt.Println("     csv2xlsx")
	fmt.Println("       -Convert all csv files in the folder to xlsx.")
	fmt.Println("     csv2xlsx ./test.csv")
	fmt.Println("       -Convert test.csv to xlsx.")
	fmt.Println()
	fmt.Println("   Ver: " + version + ", Author: NadOo")
	fmt.Println("   Source: https://github.com/nadoo/csv2xlsx")
	fmt.Println("-")
}

func main() {
	showHelp()

	if len(os.Args) > 1 {
		csv2xlsx(os.Args[1])
		return
	}

	csvFiles, _ := filepath.Glob("./*.csv")
	for _, csvFile := range csvFiles {
		csv2xlsx(csvFile)
	}
}
