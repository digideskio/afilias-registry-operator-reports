package server

import (
	"encoding/json"
	"github.com/dothiv/afilias-registry-operator-reports/repository"
	"net/http"
	"os"
)

type TransactionEventList struct {
	JsonLDTyped
	Items []*TransactionEvent `json:"items"`
	Total int                 `json:"total"`
}

type TransactionEvent struct {
	JsonLDTyped
	TLD             string
	RegistrarExtID  string
	RegistrarName   string
	ServerTrID      string
	Command         string
	ObjectType      string
	ObjectName      string
	TransactionDate string
}

type TransactionController struct {
	repo repository.TransactionRepositoryInterface
}

func (c *TransactionController) listingHandler(w http.ResponseWriter, r *http.Request) {
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

	list := new(TransactionEventList)
	list.Total = total
	list.JsonLDContext = "http://jsonld.click4life.hiv/List"
	list.JsonLDType = "http://jsonld.click4life.hiv/Afilias/TransactionEvent"
	list.JsonLDId = "/registrations"
	list.Items = make([]*TransactionEvent, len(events))

	for i := range events {
		e := new(TransactionEvent)
		e.JsonLDContext = "http://jsonld.click4life.hiv/Afilias/TransactionEvent"
		e.TLD = events[i].TLD
		e.RegistrarExtID = events[i].Registrar_Ext_ID
		e.RegistrarName = events[i].Registrar_Name
		e.ServerTrID = events[i].Server_TrID
		e.Command = events[i].Command
		e.ObjectType = events[i].Object_Type
		e.ObjectName = events[i].Object_Name
		e.TransactionDate = events[i].Transaction_Date
		list.Items[i] = e
	}

	w.Header().Add("Content-Type", "application/json")
	// Add nwext link
	if len(events) > 0 {
		last := list.Items[len(events)-1]
		w.Header().Add("Link", `</transactions?offsetKey=`+last.ServerTrID+`>; rel="next"`)
	} else {
		w.Header().Add("Link", `</transactions?offsetKey=`+maxKey+`>; rel="next"`)
	}
	encoder := json.NewEncoder(w)
	encoder.Encode(list)
}
