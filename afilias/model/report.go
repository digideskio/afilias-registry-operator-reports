package model

import (
	"fmt"
	regexputil "github.com/dothiv/afilias-registry-operator-reports/util/regexp"
	"regexp"
)

// Represents an Afilias Report definition
//
// The filename of the reports generally use the following naming convention:
//
// +--------------------+------------------------------------------------------+
// | File Name Part     | Details                                              |
// +--------------------+------------------------------------------------------+
// | <Prefix>           | All registry operator reports are represented by the |
// |                    | value of “RO”. All registrar reports are represented |
// |                    | by the value of “R”.                                 |
// +--------------------+------------------------------------------------------+
// | <Content>          | Descriptive name indicating the report content.      |
// +--------------------+------------------------------------------------------+
// | <Grouping>         | Grouping could be:                                   |
// |                    | “ALL” – Includes all TLDs sponsored by the operator  |
// |                    | <Group Name/Registry Operator Name> – Based on the   |
// |                    | name of the group of TLDs or the registry operator   |
// |                    | name                                                 |
// +--------------------+------------------------------------------------------+
// | <Frequency>        | Daily, Weekly, Monthly, Quarterly                    |
// +--------------------+------------------------------------------------------+
// | <Date>             | The date and range of the data contained in the      |
// |                    | report.                                              |
// |                    | If <Frequency> is ‘daily’ or ‘weekly’,               |
// |                    | Then format for <Date> is YYYY-MM-DD                 |
// |                    | If <Frequency> is ‘monthly’,                         |
// |                    | Then format for <Date> is YYYY-MM                    |
// |                    | If <Frequency> is ‘quarterly’,                       |
// |                    | Then format for is YYYY-Q[1,2,3,4]                   |
// +--------------------+------------------------------------------------------+
// | General Convention | <Prefix>_<Content>_<Grouping>_<Frequency>_<Date>     |
// +--------------------+------------------------------------------------------+
type Report struct {
	Prefix    string
	Content   string
	Grouping  string
	Frequency string
	Date      string
}

var PREFIX_FORMAT = regexp.MustCompile(`^RO*$`)
var FREQUENCY_FORMAT = regexp.MustCompile(`^(daily|weekly|monthly|quarterly)$`)
var DATE_FORMAT_DAILY = regexp.MustCompile(`^2[0-9]{3}-(0[1-9]|1[0-2])-(0[1-9]|[1-2][0-9]|3[0-1])$`)
var DATE_FORMAT_MONTHLY = regexp.MustCompile(`^2[0-9]{3}-(0[1-9]|1[0-2])$`)
var DATE_FORMAT_QUARTERLY = regexp.MustCompile(`^2[0-9]{3}-Q[1-4]$`)

func NewReport(prefix string, content string, grouping string, frequency string, date string) (r *Report, err error) {
	if !PREFIX_FORMAT.MatchString(prefix) {
		err = fmt.Errorf("Invalid prefix: '%s'", prefix)
		return
	}
	if !FREQUENCY_FORMAT.MatchString(frequency) {
		err = fmt.Errorf("Invalid frequency: '%s'", frequency)
		return
	}
	if frequency == "daily" || frequency == "weekly" {
		if !DATE_FORMAT_DAILY.MatchString(date) {
			err = fmt.Errorf("Invalid date: '%s'", date)
			return
		}
	}
	if frequency == "monthly" {
		if !DATE_FORMAT_MONTHLY.MatchString(date) {
			err = fmt.Errorf("Invalid date: '%s'", date)
			return
		}
	}
	if frequency == "quarterly" {
		if !DATE_FORMAT_QUARTERLY.MatchString(date) {
			err = fmt.Errorf("Invalid date: '%s'", date)
			return
		}
	}
	r = new(Report)
	r.Prefix = prefix
	r.Content = content
	r.Grouping = grouping
	r.Frequency = frequency
	r.Date = date
	return
}

func (r *Report) GetName() (name string) {
	return fmt.Sprintf("%s_%s_%s_%s_%s.txt", r.Prefix, r.Content, r.Grouping, r.Frequency, r.Date)
}

var NAME_FORMAT = regexp.MustCompile(`^` +
	"(?P<prefix>" + PREFIX_FORMAT.String()[1:len(PREFIX_FORMAT.String())-1] + ")_" +
	"(?P<content>[^_]+)_" +
	"(?P<grouping>[^_]+)_" +
	"(?P<frequency>" + FREQUENCY_FORMAT.String()[1:len(FREQUENCY_FORMAT.String())-1] + ")_" +
	"(?P<date>" +
	"(" + DATE_FORMAT_DAILY.String()[1:len(DATE_FORMAT_DAILY.String())-1] + ")" +
	"|" +
	"(" + DATE_FORMAT_MONTHLY.String()[1:len(DATE_FORMAT_MONTHLY.String())-1] + ")" +
	"|" +
	"(" + DATE_FORMAT_QUARTERLY.String()[1:len(DATE_FORMAT_QUARTERLY.String())-1] + ")" +
	")" +
	`.txt$`)

// Create a report from a filename
func NewReportFromName(name string) (r *Report, err error) {
	if !NAME_FORMAT.MatchString(name) {
		err = fmt.Errorf("Invalid name: '%s'", name)
		return
	}
	matches := regexputil.MapRegexpGroups(NAME_FORMAT, name)
	r = new(Report)
	r.Prefix = matches["prefix"]
	r.Content = matches["content"]
	r.Grouping = matches["grouping"]
	r.Frequency = matches["frequency"]
	r.Date = matches["date"]
	return
}
