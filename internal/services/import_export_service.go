package services

import (
	"fmt"
	"io"
	"reflect"

	"github.com/xuri/excelize/v2"
)

type ImportExportService interface {
	ExportToExcel(data interface{}, sheetName string) ([]byte, error)
	ReadExcel(file io.Reader, sheetName string) ([][]string, error)
}

type importExportService struct{}

func NewImportExportService() ImportExportService {
	return &importExportService{}
}

func (s *importExportService) ExportToExcel(data interface{}, sheetName string) ([]byte, error) {
	f := excelize.NewFile()
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, err
	}
	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")

	val := reflect.ValueOf(data)
	if val.Kind() != reflect.Slice {
		return nil, fmt.Errorf("data must be a slice")
	}

	if val.Len() == 0 {
		buf, err := f.WriteToBuffer()
		if err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	}

	// Get headers from struct tags or field names
	firstElem := val.Index(0)
	if firstElem.Kind() == reflect.Ptr {
		firstElem = firstElem.Elem()
	}

	if firstElem.Kind() != reflect.Struct {
		// Basic slice of maps or other types not supported yet for auto-headers
		return nil, fmt.Errorf("slice elements must be structs for auto-header export")
	}

	typ := firstElem.Type()
	headers := []string{}
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		header := field.Tag.Get("excel")
		if header == "" {
			header = field.Name
		}
		if header != "-" {
			headers = append(headers, header)
		}
	}

	// Write headers
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, header)
	}

	// Write data
	for i := 0; i < val.Len(); i++ {
		elem := val.Index(i)
		if elem.Kind() == reflect.Ptr {
			elem = elem.Elem()
		}

		colIdx := 1
		for j := 0; j < typ.NumField(); j++ {
			field := typ.Field(j)
			if field.Tag.Get("excel") == "-" {
				continue
			}

			cell, _ := excelize.CoordinatesToCellName(colIdx, i+2)
			f.SetCellValue(sheetName, cell, elem.Field(j).Interface())
			colIdx++
		}
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *importExportService) ReadExcel(file io.Reader, sheetName string) ([][]string, error) {
	f, err := excelize.OpenReader(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if sheetName == "" {
		sheetName = f.GetSheetName(0)
	}

	return f.GetRows(sheetName)
}
