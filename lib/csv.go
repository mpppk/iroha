package lib

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/pkg/errors"
)

func ReadCSVAsSlice(filename string) ([]string, error) {
	records, err := ReadCSV(filename)
	if err != nil {
		return nil, err
	}
	var words []string
	for _, record := range records {
		words = append(words, record...)
	}
	return words, nil
}

func ReadCSV(filename string) ([][]string, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read csv from %s", filename)
	}
	r := csv.NewReader(strings.NewReader(string(bytes)))

	records, err := r.ReadAll()
	return records, errors.Wrapf(err, "failed to read csv contents from %s", filename)
}
