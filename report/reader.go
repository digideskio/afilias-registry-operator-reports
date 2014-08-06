package report

import (
	"encoding/csv"
	model "github.com/dothiv/afilias-registry-operator-reports/afilias/model"
	"io"
	"os"
	"strings"
)

// Can read Afilias Reports
//
// MinFields: Ignore lines with less that this number of fields
// SkipErrors: Ignore lines which could not be parsed
type Reader struct {
	Report     *model.Report
	file       *os.File
	csvReader  *csv.Reader
	Header     []string
	MinFields  int
	SkipErrors bool
}

// Create a new reader for report files
func NewReportReader() (reader *Reader) {
	reader = new(Reader)
	reader.MinFields = 2
	reader.SkipErrors = false
	return
}

// Close open files
func (r *Reader) Close() {
	r.file.Close()
}

// Open a report file
func (r *Reader) Open(dir string, report *model.Report) (err error) {
	r.file, err = os.Open(dir + report.GetName())
	if err != nil {
		return
	}
	r.Report = report

	r.csvReader = csv.NewReader(r.file)
	r.csvReader.Comma = '|'
	if r.SkipErrors {
		r.csvReader.FieldsPerRecord = -1
	}
	r.Header, err = r.line()
	if err != nil {
		return
	}
	return
}

// Read the next line from the report
func (r *Reader) line() (line []string, err error) {
	for {
		line, err = r.csvReader.Read()
		if err == nil {
			if len(line) < r.MinFields {
				continue
			}
			return
		}
		if err == io.EOF {
			return
		}
		if r.SkipErrors {
			err = nil
			continue
		}
		return
	}
	return
}

// Read the next line from the report and maps it to the header fields
func (r *Reader) Next() (mappedLine map[string]string, err error) {
	var line []string
	line, err = r.line()
	if err != nil {
		return
	}
	mappedLine = make(map[string]string)
	for k := range r.Header {
		mappedLine[strings.TrimSpace(r.Header[k])] = strings.TrimSpace(line[k])
	}
	return
}
