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

type RegistrationEventList struct {
	JsonLDTyped
	Items []*RegistrationEvent `json:"items"`
	Total int                  `json:"total"`
}

type RegistrationEvent struct {
	JsonLDTyped
	DomainId           string
	DomainName         string
	DomainCreatedOn    string
	RegistrarExtId     string
	RegistrantClientId string
	RegistrantName     string
	RegistrantOrg      string
	RegistrantEmail    string
}

type Controller struct {
	repo *repository.Repository
}

func (c *Controller) entryPointHandler(w http.ResponseWriter, r *http.Request) {
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

func (c *Controller) registrationsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(400)
		return
	}
	formErr := r.ParseForm()
	if formErr != nil {
		w.WriteHeader(500)
		os.Stderr.WriteString(formErr.Error() + "\n")
		return
	}

	itemsPerPage := 100
	offsetKey := r.Form.Get("offsetKey")
	events, findErr := c.repo.FindPaginated(itemsPerPage, offsetKey)
	if findErr != nil {
		w.WriteHeader(500)
		os.Stderr.WriteString(findErr.Error() + "\n")
		return
	}

	total, maxKey, statsErr := c.repo.Stats()
	if statsErr != nil {
		w.WriteHeader(500)
		os.Stderr.WriteString(statsErr.Error() + "\n")
		return
	}

	list := new(RegistrationEventList)
	list.Total = total
	list.JsonLDContext = "http://jsonld.click4life.hiv/List"
	list.JsonLDType = "http://jsonld.click4life.hiv/Afilias/RegistrationEvent"
	list.JsonLDId = "/registrations"
	list.Items = make([]*RegistrationEvent, len(events))

	for i := range events {
		e := new(RegistrationEvent)
		e.JsonLDContext = "http://jsonld.click4life.hiv/Afilias/RegistrationEvent"
		e.DomainId = events[i].DomainId
		e.DomainName = events[i].DomainName
		e.DomainCreatedOn = events[i].DomainCreatedOn
		e.RegistrarExtId = events[i].RegistrarExtId
		e.RegistrantClientId = events[i].RegistrantClientId
		e.RegistrantName = events[i].RegistrantName
		e.RegistrantOrg = events[i].RegistrantOrg
		e.RegistrantEmail = events[i].RegistrantEmail
		list.Items[i] = e
	}

	w.Header().Add("Content-Type", "application/json")
	// Add nwext link
	if len(events) > 0 {
		last := list.Items[len(events)-1]
		w.Header().Add("Link", `</registrations?offsetKey=`+last.DomainId+`>; rel="next"`)
	} else {
		w.Header().Add("Link", `</registrations?offsetKey=`+maxKey+`>; rel="next"`)
	}
	encoder := json.NewEncoder(w)
	encoder.Encode(list)
}

func Run(c *Config) (err error) {
	// Open DB
	db, err := sql.Open("postgres", c.Database.DSN())
	if err != nil {
		return
	}
	cntrl := new(Controller)
	cntrl.repo = repository.NewDomainContactDetailsHourlyRepository(db)

	os.Stdout.WriteString(fmt.Sprintf("Running server on localhost:%d\n", c.Port))
	http.HandleFunc("/registrations", cntrl.registrationsHandler)
	http.HandleFunc("/", cntrl.entryPointHandler)
	http.ListenAndServe(fmt.Sprintf("localhost:%d", c.Port), nil)
	return
}
