package server

import (
	"code.google.com/p/gcfg"
	"database/sql"
	"encoding/json"
	registrations "github.com/dothiv/afilias-registry-operator-reports/command/importer/registrations"
	transactions "github.com/dothiv/afilias-registry-operator-reports/command/importer/transactions"
	"github.com/dothiv/afilias-registry-operator-reports/repository"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Test for registrations

type RegistrationList struct {
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

	cntrl := new(RegistrationController)
	db, _ := sql.Open("postgres", c.Database.DSN())
	cntrl.repo = repository.NewDomainContactDetailsHourlyRepository(db)

	ts := httptest.NewServer(http.HandlerFunc(cntrl.listingHandler))
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

	var l RegistrationList
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

func TestThatItReturnsNextUrlAfterEndOfRegistrations(t *testing.T) {
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

	cntrl := new(RegistrationController)
	db, _ := sql.Open("postgres", c.Database.DSN())
	cntrl.repo = repository.NewDomainContactDetailsHourlyRepository(db)

	ts := httptest.NewServer(http.HandlerFunc(cntrl.listingHandler))
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

// Test for transactions

type TransactionList struct {
	Total int
	Items []struct {
		TLD             string
		RegistrarExtID  string
		RegistrarName   string
		ServerTrID      string
		Command         string
		ObjectType      string
		ObjectName      string
		TransactionDate string
	}
}

func TestThatItListsTransactions(t *testing.T) {
	assert := assert.New(t)

	// Import
	c := transactions.NewDefaultConfig()
	c.ConfigFile = "../../test.ini"
	c.ReportsDir = "../../example"
	c.Quiet = true
	configErr := gcfg.ReadFileInto(c, c.ConfigFile)
	if configErr != nil {
		t.Fatal(configErr)
	}

	importErr := transactions.Import(c)
	if importErr != nil {
		t.Fatal(importErr)
	}

	cntrl := new(TransactionController)
	db, _ := sql.Open("postgres", c.Database.DSN())
	cntrl.repo = repository.NewTransactionRepository(db)

	ts := httptest.NewServer(http.HandlerFunc(cntrl.listingHandler))
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

	var l TransactionList
	unmarshalErr := json.Unmarshal(b, &l)
	if unmarshalErr != nil {
		t.Fatal(unmarshalErr)
	}
	assert.Equal(9, l.Total)

	assert.Equal("hiv", l.Items[0].TLD)
	assert.Equal("1155-YN", l.Items[0].RegistrarExtID)
	assert.Equal("Whois Networks Co., Ltd.", l.Items[0].RegistrarName)
	assert.Equal("26515244", l.Items[0].ServerTrID)
	assert.Equal("CREATE", l.Items[0].Command)
	assert.Equal("DOMAIN", l.Items[0].ObjectType)
	assert.Equal("samsung.hiv", l.Items[0].ObjectName)
	assert.Equal("2014-07-25 08:04:03", l.Items[0].TransactionDate)

	assert.Equal("hiv", l.Items[5].TLD)
	assert.Equal("1180-LL", l.Items[5].RegistrarExtID)
	assert.Equal("Lexsynergy Limited", l.Items[5].RegistrarName)
	assert.Equal("26520559", l.Items[5].ServerTrID)
	assert.Equal("CREATE", l.Items[5].Command)
	assert.Equal("DOMAIN", l.Items[5].ObjectType)
	assert.Equal("stanbic.hiv", l.Items[5].ObjectName)
	assert.Equal("2014-07-25 14:42:41", l.Items[5].TransactionDate)

	assert.Equal(`</transactions?offsetKey=27568803>; rel="next"`, res.Header.Get("Link"))
}

func TestThatItReturnsNextUrlAfterEndOfTransactions(t *testing.T) {
	assert := assert.New(t)

	// Import
	c := transactions.NewDefaultConfig()
	c.ConfigFile = "../../test.ini"
	c.ReportsDir = "../../example"
	c.Quiet = true
	configErr := gcfg.ReadFileInto(c, c.ConfigFile)
	if configErr != nil {
		t.Fatal(configErr)
	}

	importErr := transactions.Import(c)
	if importErr != nil {
		t.Fatal(importErr)
	}

	cntrl := new(TransactionController)
	db, _ := sql.Open("postgres", c.Database.DSN())
	cntrl.repo = repository.NewTransactionRepository(db)

	ts := httptest.NewServer(http.HandlerFunc(cntrl.listingHandler))
	defer ts.Close()

	res, err := http.Get(ts.URL + "/registrations?offsetKey=27568803")
	if err != nil {
		t.Fatal(err)
	}
	_, err = ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(`</transactions?offsetKey=27568803>; rel="next"`, res.Header.Get("Link"))

}
