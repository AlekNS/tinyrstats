package helpers

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/alekns/tinyrstats/internal/monitor"
)

// ReadTasksFromCsvFile reads and creates tasks from csv file.
func ReadTasksFromCsvFile(filePath string) ([]*monitor.ScheduleHealthTask, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	results := make([]*monitor.ScheduleHealthTask, 0)

	csvReader := csv.NewReader(bufio.NewReader(f))
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

		results = append(results, &monitor.ScheduleHealthTask{
			Interval: 0,
			Task: &monitor.HealthTask{
				URL: row[0],
			},
		})
	}

	return results, nil
}
