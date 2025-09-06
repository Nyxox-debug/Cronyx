package loaders

import (
	"context"
	"encoding/csv"
	"github.com/Nyxox-debug/Cronyx/pkg/cronyx"
	"io"
	"os"
)

type CSVLoader struct{}

func (c CSVLoader) Load(ctx context.Context, cfg cronyx.DataSourceConfig) (cronyx.DataPayload, error) {
	path := cfg["path"]
	f, err := os.Open(path)
	if err != nil {
		return cronyx.DataPayload{}, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	headers, err := r.Read()
	if err != nil {
		return cronyx.DataPayload{}, err
	}
	var rows []map[string]interface{}
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return cronyx.DataPayload{}, err
		}
		row := map[string]interface{}{}
		for i, h := range headers {
			row[h] = record[i]
		}
		rows = append(rows, row)
	}
	return cronyx.DataPayload{Rows: rows}, nil
}
