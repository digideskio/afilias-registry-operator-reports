package transactions

import (
	"code.google.com/p/gcfg"
	"database/sql"
	"github.com/dothiv/afilias-registry-operator-reports/repository"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Integration test
func TestThatDataIsImported(t *testing.T) {
	assert := assert.New(t)

	c := NewDefaultConfig()
	c.ConfigFile = "../../../test.ini"
	c.ReportsDir = "../../../example"
	configErr := gcfg.ReadFileInto(c, c.ConfigFile)
	assert.Nil(configErr)

	importErr := Import(c)
	assert.Nil(importErr)

	// Verify import
	db, _ := sql.Open("postgres", c.Database.DSN())
	repo := repository.NewTransactionRepository(db)
	transactions, findErr := repo.FindAll()
	assert.Nil(findErr)
	assert.Equal(9, len(transactions))

	assert.Equal("hiv", transactions[0].TLD)
	assert.Equal("1155-YN", transactions[0].Registrar_Ext_ID)
	assert.Equal("Whois Networks Co., Ltd.", transactions[0].Registrar_Name)
	assert.Equal("26515244", transactions[0].Server_TrID)
	assert.Equal("CREATE", transactions[0].Command)
	assert.Equal("DOMAIN", transactions[0].Object_Type)
	assert.Equal("samsung.hiv", transactions[0].Object_Name)
	assert.Equal("2014-07-25 08:04:03", transactions[0].Transaction_Date)

	assert.Equal("hiv", transactions[5].TLD)
	assert.Equal("1180-LL", transactions[5].Registrar_Ext_ID)
	assert.Equal("Lexsynergy Limited", transactions[5].Registrar_Name)
	assert.Equal("26520559", transactions[5].Server_TrID)
	assert.Equal("CREATE", transactions[5].Command)
	assert.Equal("DOMAIN", transactions[5].Object_Type)
	assert.Equal("stanbic.hiv", transactions[5].Object_Name)
	assert.Equal("2014-07-25 14:42:41", transactions[5].Transaction_Date)

}
