package server

import (
	"encoding/json"
	"github.com/dothiv/afilias-registry-operator-reports/repository"
	"net/http"
	"os"
)

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

type RegistrationController struct {
	repo repository.DomainContactDetailsHouryRepositoryInterface
}

func (c *RegistrationController) listingHandler(w http.ResponseWriter, r *http.Request) {
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
