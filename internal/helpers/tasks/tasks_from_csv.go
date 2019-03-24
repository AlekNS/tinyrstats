package tasks

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"

	"github.com/alekns/tinyrstats/internal/monitor"
)

// ReadTasksFromCsvFile reads and creates tasks from csv file.
func ReadTasksFromCsvFile(defaultProtocol, filePath string) ([]*monitor.ScheduleHealthTask, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	results := make([]*monitor.ScheduleHealthTask, 0)
	csvReader := csv.NewReader(bufio.NewReader(f))
	csvReader.Comma = ';'
	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if len(row) > 1 {
			return nil, fmt.Errorf("invalid csv format, only single column supported yet")
		}

		if len(row) == 0 {
			continue
		}

		parsedURL, err := url.Parse(row[0])
		if err != nil {
			return nil, fmt.Errorf("invalid resource format: %v: %v", row[0], err)
		}

		if len(parsedURL.Scheme) == 0 {
			parsedURL.Scheme = defaultProtocol
		}

		if !strings.HasPrefix(parsedURL.String(), "http://") &&
			!strings.HasPrefix(parsedURL.String(), "https://") {
			return nil, fmt.Errorf("invalid resource row: %v", row[0])
		}

		results = append(results, &monitor.ScheduleHealthTask{
			Interval: 0,
			Task: &monitor.HealthTask{
				URL:    parsedURL.String(),
				Method: "GET",
			},
		})
	}

	return results, nil
}
