package registrations

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

	c := newDefaultConfig()
	c.ConfigFile = "../../../test.ini"
	c.ReportsDir = "../../../example"
	configErr := gcfg.ReadFileInto(c, c.ConfigFile)
	assert.Nil(configErr)

	importErr := Import(c)
	assert.Nil(importErr)

	// Verify import
	db, _ := sql.Open("postgres", c.Database.DSN())
	repo := repository.NewDomainContactDetailsHourlyRepository(db)
	events, findErr := repo.FindAll()
	assert.Nil(findErr)
	assert.Equal(2, len(events))

	assert.Equal("2863429", events[0].DomainId)
	assert.Equal("bcme.hiv", events[0].DomainName)
	assert.Equal("2014-07-21 17:34:18.349+00", events[0].DomainCreatedOn)
	assert.Equal("1061-EM", events[0].RegistrarExtId)
	assert.Equal("mmr-138842", events[0].RegistrantClientId)
	assert.Equal("Domain Administrator", events[0].RegistrantName)
	assert.Equal("Bcme LLC.", events[0].RegistrantOrg)
	assert.Equal("domain@bcme.com", events[0].RegistrantEmail)

	assert.Equal("2863499", events[1].DomainId)
	assert.Equal("acme.hiv", events[1].DomainName)
	assert.Equal("2014-07-21 17:43:29.824+00", events[1].DomainCreatedOn)
	assert.Equal("1061-EM", events[1].RegistrarExtId)
	assert.Equal("mmr-105291", events[1].RegistrantClientId)
	assert.Equal("Domain Administrator", events[1].RegistrantName)
	assert.Equal("Acme Inc.", events[1].RegistrantOrg)
	assert.Equal("ccops@acme.com", events[1].RegistrantEmail)
}
