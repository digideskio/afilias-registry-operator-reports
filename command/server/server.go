package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/dothiv/afilias-registry-operator-reports/repository"
	"net/http"
	"os"
)

type JsonLDTyped struct {
	JsonLDContext string `json:"@context,omitempty"`
	JsonLDId      string `json:"@id,omitempty"`
	JsonLDType    string `json:"@type,omitempty"`
}

type EntryPoint struct {
	JsonLDContext string       `json:"@context"`
	Registrations *JsonLDTyped `json:"registrations"`
}

type EntryPointController struct {
}

func (c *EntryPointController) entryPointHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(400)
		return
	}
	entryPoint := new(EntryPoint)
	entryPoint.JsonLDContext = "http://jsonld.click4life.hiv/EntryPoint"
	entryPoint.Registrations = new(JsonLDTyped)
	entryPoint.Registrations.JsonLDContext = "http://jsonld.click4life.hiv/List"
	entryPoint.Registrations.JsonLDType = "http://jsonld.click4life.hiv/Afilias/RegistrationEvent"
	entryPoint.Registrations.JsonLDId = "/registrations"
	w.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.Encode(entryPoint)
}

func Run(c *Config) (err error) {
	// Open DB
	db, err := sql.Open("postgres", c.Database.DSN())
	if err != nil {
		return
	}

	os.Stdout.WriteString(fmt.Sprintf("Running server on localhost:%d\n", c.Port))

	registrationsCntrl := new(RegistrationController)
	registrationsCntrl.repo = repository.NewDomainContactDetailsHourlyRepository(db)
	http.HandleFunc("/registrations", registrationsCntrl.listingHandler)

	transactionCntrl := new(TransactionController)
	transactionCntrl.repo = repository.NewTransactionRepository(db)
	http.HandleFunc("/transactions", transactionCntrl.listingHandler)

	entryPointCntrl := new(EntryPointController)
	http.HandleFunc("/", entryPointCntrl.entryPointHandler)

	http.ListenAndServe(fmt.Sprintf("localhost:%d", c.Port), nil)
	return
}
