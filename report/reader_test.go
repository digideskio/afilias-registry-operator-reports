package report

import (
	model "github.com/dothiv/afilias-registry-operator-reports/afilias/model"
	assert "github.com/stretchr/testify/assert"
	"io"
	"testing"
)

// Test for a regular report
func TestRead(t *testing.T) {
	assert := assert.New(t)
	report, reportErr := model.NewReport("RO", "transactions", "ALL", "daily", "2014-07-25", "txt")
	assert.Nil(reportErr)
	reader := NewReportReader()
	readerErr := reader.Open("../example/", report)
	assert.Nil(readerErr)
	defer reader.Close()
	lines := make([]map[string]string, 0)
	for {
		line, err := reader.Next()
		if err != nil {
			assert.Equal(io.EOF, err)
			break
		}
		lines = append(lines, line)
	}
	assert.Equal(7, len(lines))
	// Check header
	assert.Equal(8, len(reader.Header))
	assert.Equal("TLD", reader.Header[0])
	assert.Equal("Registrar_Ext_ID", reader.Header[1])
	assert.Equal("Registrar_Name", reader.Header[2])
	assert.Equal("Server_TrID", reader.Header[3])
	assert.Equal("Command", reader.Header[4])
	assert.Equal("Object_Type", reader.Header[5])
	assert.Equal("Object_Name", reader.Header[6])
	assert.Equal("Transaction_Date", reader.Header[7])
	// Check data
	assert.Equal(8, len(lines[0]))
	// Line 1
	assert.Equal("hiv", lines[0]["TLD"])
	assert.Equal("1155-YN", lines[0]["Registrar_Ext_ID"])
	assert.Equal("Whois Networks Co., Ltd.", lines[0]["Registrar_Name"])
	assert.Equal("26515244", lines[0]["Server_TrID"])
	assert.Equal("CREATE", lines[0]["Command"])
	assert.Equal("DOMAIN", lines[0]["Object_Type"])
	assert.Equal("samsung.hiv", lines[0]["Object_Name"])
	assert.Equal("2014-07-25 08:04:03", lines[0]["Transaction_Date"])
	// Line 7
	assert.Equal("hiv", lines[6]["TLD"])
	assert.Equal("1180-LL", lines[6]["Registrar_Ext_ID"])
	assert.Equal("Lexsynergy Limited", lines[6]["Registrar_Name"])
	assert.Equal("26520559", lines[6]["Server_TrID"])
	assert.Equal("CREATE", lines[6]["Command"])
	assert.Equal("DOMAIN", lines[6]["Object_Type"])
	assert.Equal("stanbic.hiv", lines[6]["Object_Name"])
	assert.Equal("2014-07-25 14:42:41", lines[6]["Transaction_Date"])
}

// Test for an hourly report
func TestReadHourly(t *testing.T) {
	assert := assert.New(t)
	report, reportErr := model.NewReport("RO", "domain-contact-details", "hiv", "hourly", "2014-07-21T18", "csv")
	assert.Nil(reportErr)
	reader := NewReportReader()
	readerErr := reader.Open("../example/", report)
	assert.Nil(readerErr)
	defer reader.Close()
	lines := make([]map[string]string, 0)
	for {
		line, err := reader.Next()
		if err != nil {
			assert.Equal(io.EOF, err)
			break
		}
		lines = append(lines, line)
	}
	assert.Equal(2, len(lines))
	// Check data
	assert.Equal(12, len(lines[0]))
	// Line 1
	assert.Equal("2863429", lines[0]["domain_id"])
	assert.Equal("bcme.hiv", lines[0]["domain_name"])
	assert.Equal("2014-07-21 17:34:18.349+00", lines[0]["domain_created_on"])
	assert.Equal("1061-EM", lines[0]["registrar_ext_id"])
	assert.Equal("", lines[0]["ipr_name"])
	assert.Equal("", lines[0]["ipr_number"])
	assert.Equal("", lines[0]["ipr_class"])
	assert.Equal("mmr-138842", lines[0]["registrant_client_id"])
	assert.Equal("Domain Administrator", lines[0]["registrant_name"])
	assert.Equal("Bcme LLC.", lines[0]["registrant_org"])
	assert.Equal("domain@bcme.com", lines[0]["registrant_email"])
	assert.Equal("", lines[0]["member_id"])
	// Line 2
	assert.Equal("2863499", lines[1]["domain_id"])
	assert.Equal("acme.hiv", lines[1]["domain_name"])
	assert.Equal("2014-07-21 17:43:29.824+00", lines[1]["domain_created_on"])
	assert.Equal("1061-EM", lines[1]["registrar_ext_id"])
	assert.Equal("", lines[1]["ipr_name"])
	assert.Equal("", lines[1]["ipr_number"])
	assert.Equal("", lines[1]["ipr_class"])
	assert.Equal("mmr-105291", lines[1]["registrant_client_id"])
	assert.Equal("Domain Administrator", lines[1]["registrant_name"])
	assert.Equal("Acme Inc.", lines[1]["registrant_org"])
	assert.Equal("ccops@acme.com", lines[1]["registrant_email"])
	assert.Equal("", lines[1]["member_id"])
}

// Test for an deposit report
func TestReadDeposit(t *testing.T) {
	assert := assert.New(t)
	report, reportErr := model.NewReport("RO", "deposit-withdrawal", "hiv", "monthly", "2014-08", "txt")
	assert.Nil(reportErr)
	reader := NewReportReader()
	reader.SkipErrors = true
	readerErr := reader.Open("../example/", report)
	assert.Nil(readerErr)
	defer reader.Close()
	reader.SkipErrors = true
	lines := make([]map[string]string, 0)
	for {
		line, err := reader.Next()
		if err != nil {
			assert.Equal(io.EOF, err)
			break
		}
		lines = append(lines, line)
	}
	assert.Equal(43, len(lines))

	// Check data
	assert.Equal(6, len(lines[0]))
	// Line 1
	assert.Equal("Admin Contact", lines[0]["Done By"])
	assert.Equal("1001-HV", lines[0]["Reg ID"])
	assert.Equal("TLD dotHIV Registry GmbH", lines[0]["Registrar"])
	assert.Equal("2014-07-15 22:27:49.05003", lines[0]["Transaction Date"])
	assert.Equal("initial deposit", lines[0]["Type"])
	assert.Equal("1000000.00", lines[0]["Amount"])
}
