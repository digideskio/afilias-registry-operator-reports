package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestThatICanCreateANewReport(t *testing.T) {
	assert := assert.New(t)
	r, err := NewReport("RO", "transaction", "ALL", "daily", "2013-10-31")
	assert.Nil(err)
	assert.Equal("RO_transaction_ALL_daily_2013-10-31.txt", r.GetName())
}

func TestThatICanNotCreateAReportWithInvalidPrefix(t *testing.T) {
	assert := assert.New(t)
	r, err := NewReport("A", "transaction", "ALL", "daily", "2013-10-31")
	assert.NotNil(err)
	assert.Nil(r)
	assert.Equal("Invalid prefix: 'A'", err.Error())
}

func TestThatICanNotCreateAReportWithInvalidFrequency(t *testing.T) {
	assert := assert.New(t)
	r, err := NewReport("RO", "transaction", "ALL", "bi-weekly", "2013-10-31")
	assert.NotNil(err)
	assert.Equal("Invalid frequency: 'bi-weekly'", err.Error())
	assert.Nil(r)
}

func TestThatICanNotCreateAReportWithInvalidDate(t *testing.T) {
	assert := assert.New(t)
	invalidDates := [][]string{[]string{"monthly", "2013-10-31"}, []string{"daily", "2013-10"}, []string{"quarterly", "2013-10-31"}}
	for i := range invalidDates {
		r, err := NewReport("RO", "transaction", "ALL", invalidDates[i][0], invalidDates[i][1])
		assert.NotNil(err)
		assert.Equal("Invalid date: '"+invalidDates[i][1]+"'", err.Error())
		assert.Nil(r)
	}
}

func TestThatICanCreateAReportWithValidDate(t *testing.T) {
	assert := assert.New(t)
	invalidDates := [][]string{[]string{"weekly", "2013-10-31"}, []string{"monthly", "2013-10"}, []string{"quarterly", "2013-Q2"}}
	for i := range invalidDates {
		r, err := NewReport("RO", "transaction", "ALL", invalidDates[i][0], invalidDates[i][1])
		assert.Nil(err)
		assert.NotNil(r)
	}
}

func TestThatICanCreateANewReportFromAName(t *testing.T) {
	assert := assert.New(t)
	r, err := NewReportFromName("RO_transaction_ALL_daily_2013-10-31.txt")
	assert.Nil(err)
	assert.Equal("RO", r.Prefix)
	assert.Equal("transaction", r.Content)
	assert.Equal("ALL", r.Grouping)
	assert.Equal("daily", r.Frequency)
	assert.Equal("2013-10-31", r.Date)
}
