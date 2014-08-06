package report

import (
	"encoding/csv"
	model "github.com/dothiv/afilias-registry-operator-reports/afilias/model"
	"os"
)

// Can read Afilias Reports
type Reader struct {
	Report    *model.Report
	file      *os.File
	csvReader *csv.Reader
	Header    []string
}

// Create a new reader for report files
func NewReportReader(dir string, report *model.Report) (reader *Reader, err error) {
	reader = new(Reader)
	reader.file, err = os.Open(dir + report.GetName())
	if err != nil {
		return
	}
	reader.Report = report
	reader.Header, err = reader.line()
	if err != nil {
		return
	}
	return
}

// Returns (and inits) the csv reader for the file
func (r *Reader) getCsvReader() *csv.Reader {
	if r.csvReader == nil {
		r.csvReader = csv.NewReader(r.file)
		r.csvReader.Comma = '|'
	}
	return r.csvReader
}

// Read the next line from the report
func (r *Reader) line() (line []string, err error) {
	csvReader := r.getCsvReader()
	line, err = csvReader.Read()
	if err != nil {
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
		mappedLine[r.Header[k]] = line[k]
	}
	return
}

// Close open files
func (r *Reader) Close() {
	r.file.Close()
}
