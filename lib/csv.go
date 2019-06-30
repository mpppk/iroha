package lib

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"strings"
)

func ReadCSV(filename string) ([]string, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read csv from %s", filename)
	}
	r := csv.NewReader(strings.NewReader(string(bytes)))

	records, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read csv contents from %s", filename)
	}
	var words []string
	for _, record := range records {
		words = append(words, record...)
	}
	return words, nil
}