package model

type DomainContactDetailsHourly struct {
	Model
	DomainId           string
	DomainName         string
	DomainCreatedOn    string
	RegistrarExtId     string
	RegistrantClientId string
	RegistrantName     string
	RegistrantOrg      string
	RegistrantEmail    string
}
