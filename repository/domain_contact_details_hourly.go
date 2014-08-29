package repository

import (
	"database/sql"
	"github.com/dothiv/afilias-registry-operator-reports/afilias/model"
	_ "github.com/lib/pq"
)

const TABLE_NAME = "domain_contact_details_hourly"

type DomainContactDetailsHourlyRepository struct {
	Repository
}

func NewDomainContactDetailsHourlyRepository(db *sql.DB) (repo *DomainContactDetailsHourlyRepository) {
	repo = new(DomainContactDetailsHourlyRepository)
	repo.db = db
	return
}

func (repo *DomainContactDetailsHourlyRepository) GetLatestImportTime() (domain_created_on string, err error) {
	err = repo.db.QueryRow("SELECT domain_created_on FROM " + TABLE_NAME + " ORDER BY domain_created_on DESC LIMIT 1").Scan(&domain_created_on)
	return
}

func (repo *DomainContactDetailsHourlyRepository) Persist(m *model.DomainContactDetailsHourly) (result sql.Result, err error) {
	result, err = repo.db.Exec("INSERT INTO "+TABLE_NAME+" "+
		"(domain_id, domain_name, domain_created_on, registrar_ext_id, registrant_client_id, registrant_name, registrant_org, registrant_email) "+
		"VALUES($1, $2, $3, $4, $5, $6, $7, $8)",
		m.DomainId,
		m.DomainName,
		m.DomainCreatedOn,
		m.RegistrarExtId,
		m.RegistrantClientId,
		m.RegistrantName,
		m.RegistrantOrg,
		m.RegistrantEmail)
	return
}
