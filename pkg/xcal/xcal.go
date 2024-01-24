package xcal

import (
	"bytes"
	"encoding/csv"
	"strconv"

	"github.com/elliotchance/orderedmap/v2"
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
	data := orderedmap.NewOrderedMap[string, []string]()
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
			if v, ok := data.Get(sheet.Name); ok {
				data.Set(sheet.Name, append(v, *area))
			} else {
				data.Set(sheet.Name, []string{*area})
			}
		}
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	defer writer.Flush()

	headerRow := []string{"Sample"}
	for el := data.Front(); el != nil; el = el.Next() {
		headerRow = append(headerRow, el.Key)
	}
	if err = writer.Write(headerRow); err != nil {
		return nil, errors.Wrap(err, "cannot write header row")
	}

	for i := 0; i < len(samples); i++ {
		dataRow := []string{samples[i]}
		for el := data.Front(); el != nil; el = el.Next() {
			dataRow = append(dataRow, el.Value[i])
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
