package report

import (
	model "github.com/dothiv/afilias-registry-operator-reports/afilias/model"
	assert "github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func TestRead(t *testing.T) {
	assert := assert.New(t)
	report, _ := model.NewReport("RO", "transactions", "ALL", "daily", "2014-07-25")
	reader, readerErr := NewReportReader("../example/", report)
	assert.Nil(readerErr)
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
	// Line 1
	assert.Equal(8, len(lines[0]))
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
