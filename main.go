package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
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
	xlsx.SetDefaultFont(10, "Verdana")

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

	if err != nil && err != io.EOF {
		fmt.Printf(err.Error())
		return
	}

	fileName := strings.TrimSuffix(filepath.Base(csvPath), ".csv")
	outFile := fileName + ".xlsx"
	xlsxFile.Save(outFile)
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
	fmt.Println("   Ver: 1.0, Author: NadOo")
	fmt.Println("   Source: https://github.com/nadoo/csv2xlsx")
	fmt.Println("-")
	fmt.Println()
}

func main() {
	showHelp()

	if len(os.Args) > 1 {
		csv2xlsx(os.Args[1])
		return
	}

	csvFiles, _ := filepath.Glob("./*.csv")
	for _, csvFile := range csvFiles {
		fmt.Printf("Processing csv file: %s\n", csvFile)
		csv2xlsx(csvFile)
	}
}
