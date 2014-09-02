package server

import (
	"code.google.com/p/gcfg"
	"database/sql"
	"encoding/json"
	"github.com/dothiv/afilias-registry-operator-reports/command/importer/registrations"
	"github.com/dothiv/afilias-registry-operator-reports/repository"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type List struct {
	Total int
	Items []struct {
		DomainId           string
		DomainName         string
		DomainCreatedOn    string
		RegistrarExtId     string
		RegistrantClientId string
		RegistrantName     string
		RegistrantOrg      string
		RegistrantEmail    string
	}
}

func TestThatItListsRegistrations(t *testing.T) {
	assert := assert.New(t)

	// Import
	c := registrations.NewDefaultConfig()
	c.ConfigFile = "../../test.ini"
	c.ReportsDir = "../../example"
	c.Quiet = true
	configErr := gcfg.ReadFileInto(c, c.ConfigFile)
	if configErr != nil {
		t.Fatal(configErr)
	}

	importErr := registrations.Import(c)
	if importErr != nil {
		t.Fatal(importErr)
	}

	cntrl := new(Controller)
	db, _ := sql.Open("postgres", c.Database.DSN())
	cntrl.repo = repository.NewDomainContactDetailsHourlyRepository(db)

	ts := httptest.NewServer(http.HandlerFunc(cntrl.registrationsHandler))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	b, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal("application/json", res.Header.Get("Content-Type"))

	var l List
	unmarshalErr := json.Unmarshal(b, &l)
	if unmarshalErr != nil {
		t.Fatal(unmarshalErr)
	}
	assert.Equal(2, l.Total)

	assert.Equal("2863429", l.Items[0].DomainId)
	assert.Equal("bcme.hiv", l.Items[0].DomainName)
	assert.Equal("2014-07-21 17:34:18.349+00", l.Items[0].DomainCreatedOn)
	assert.Equal("1061-EM", l.Items[0].RegistrarExtId)
	assert.Equal("mmr-138842", l.Items[0].RegistrantClientId)
	assert.Equal("Domain Administrator", l.Items[0].RegistrantName)
	assert.Equal("Bcme LLC.", l.Items[0].RegistrantOrg)
	assert.Equal("domain@bcme.com", l.Items[0].RegistrantEmail)

	assert.Equal("2863499", l.Items[1].DomainId)
	assert.Equal("acme.hiv", l.Items[1].DomainName)
	assert.Equal("2014-07-21 17:43:29.824+00", l.Items[1].DomainCreatedOn)
	assert.Equal("1061-EM", l.Items[1].RegistrarExtId)
	assert.Equal("mmr-105291", l.Items[1].RegistrantClientId)
	assert.Equal("Domain Administrator", l.Items[1].RegistrantName)
	assert.Equal("Acme Inc.", l.Items[1].RegistrantOrg)
	assert.Equal("ccops@acme.com", l.Items[1].RegistrantEmail)

	assert.Equal(`</registrations?offsetKey=2863499>; rel="next"`, res.Header.Get("Link"))
}

func TestThatItReturnsNextUrlAfterEnd(t *testing.T) {
	assert := assert.New(t)

	// Import
	c := registrations.NewDefaultConfig()
	c.ConfigFile = "../../test.ini"
	c.ReportsDir = "../../example"
	c.Quiet = true
	configErr := gcfg.ReadFileInto(c, c.ConfigFile)
	if configErr != nil {
		t.Fatal(configErr)
	}

	importErr := registrations.Import(c)
	if importErr != nil {
		t.Fatal(importErr)
	}

	cntrl := new(Controller)
	db, _ := sql.Open("postgres", c.Database.DSN())
	cntrl.repo = repository.NewDomainContactDetailsHourlyRepository(db)

	ts := httptest.NewServer(http.HandlerFunc(cntrl.registrationsHandler))
	defer ts.Close()

	res, err := http.Get(ts.URL + "/registrations?offsetKey=2863499")
	if err != nil {
		t.Fatal(err)
	}
	_, err = ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(`</registrations?offsetKey=2863499>; rel="next"`, res.Header.Get("Link"))

}
