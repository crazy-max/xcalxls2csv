package xcal

import (
	"bytes"
	"encoding/csv"
	"strconv"

	"github.com/extrame/xls"
	"github.com/pkg/errors"
)

const (
	colSample = 0
	colArea   = 4
)

func ConvertToCSV(xcalxls string) ([]byte, error) {
	f, err := xls.Open(xcalxls, "utf-8")
	if err != nil {
		return nil, errors.Wrapf(err, "cannot open Xcalibur XLS file %q", xcalxls)
	}

	var samples []string
	var data = make(map[string][]string)
	for s := 0; s <= f.NumSheets()-3; s++ {
		sheet := f.GetSheet(s)
		if sheet == nil {
			continue
		}
		for r := 5; r <= int(sheet.MaxRow)-6; r++ {
			row := sheet.Row(r)
			if row == nil {
				continue
			}
			sample := getSample(row)
			if sample == nil {
				continue
			}
			area, err := getArea(row)
			if err != nil {
				return nil, errors.Wrapf(err, "cannot get area for row %d", r)
			} else if area == nil {
				continue
			}
			if s == 0 {
				samples = append(samples, *sample)
			}
			data[sheet.Name] = append(data[sheet.Name], *area)
		}
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	defer writer.Flush()

	hrow := []string{"Sample"}
	for h := range data {
		hrow = append(hrow, h)
	}
	if err = writer.Write(hrow); err != nil {
		return nil, errors.Wrap(err, "cannot write header row")
	}

	for i := 0; i < len(samples); i++ {
		dataRow := []string{samples[i]}
		for _, val := range data {
			dataRow = append(dataRow, val[i])
		}
		if err = writer.Write(dataRow); err != nil {
			return nil, errors.Wrapf(err, "cannot write data row %d", i)
		}
	}

	return buf.Bytes(), nil
}

func getSample(row *xls.Row) *string {
	sample := row.Col(colSample)
	if sample == "" {
		return nil
	}
	return &sample
}

func getArea(row *xls.Row) (*string, error) {
	area := row.Col(colArea)
	if area == "" {
		return nil, nil
	}
	if area == "NF" {
		area = "0"
	}
	if _, err := strconv.ParseFloat(area, 64); err != nil {
		return nil, errors.Wrapf(err, "cannot parse area %q", area)
	}
	return &area, nil
}
