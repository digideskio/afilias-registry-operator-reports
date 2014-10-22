package repository

import (
	"database/sql"
	"github.com/dothiv/afilias-registry-operator-reports/afilias/model"
	_ "github.com/lib/pq"
)

type DomainContactDetailsHouryRepositoryInterface interface {
	GetLatestImportTime() (domain_created_on string, err error)
	Persist(m *model.DomainContactDetailsHourly) (result sql.Result, err error)
	FindAll() (result []*model.DomainContactDetailsHourly, err error)
	Stats() (count int, maxKey string, err error)
	FindPaginated(numitems int, offsetKey string) (result []*model.DomainContactDetailsHourly, err error)
}

type DomainContactDetailsHouryRepository struct {
	DomainContactDetailsHouryRepositoryInterface
	db         *sql.DB
	TABLE_NAME string
	FIELDS     string
}

func NewDomainContactDetailsHourlyRepository(db *sql.DB) (repo *DomainContactDetailsHouryRepository) {
	repo = new(DomainContactDetailsHouryRepository)
	repo.db = db
	repo.TABLE_NAME = "domain_contact_details_hourly"
	repo.FIELDS = "domain_id, domain_name, domain_created_on, registrar_ext_id, registrant_client_id, registrant_name, registrant_org, registrant_email"
	return
}

func (repo *DomainContactDetailsHouryRepository) GetLatestImportTime() (domain_created_on string, err error) {
	err = repo.db.QueryRow("SELECT domain_created_on FROM " + repo.TABLE_NAME + " ORDER BY domain_created_on DESC LIMIT 1").Scan(&domain_created_on)
	return
}

func (repo *DomainContactDetailsHouryRepository) Persist(m *model.DomainContactDetailsHourly) (result sql.Result, err error) {
	result, err = repo.db.Exec("INSERT INTO "+repo.TABLE_NAME+" "+
		"("+repo.FIELDS+") "+
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

func (repo *DomainContactDetailsHouryRepository) rowsToResult(rows *sql.Rows) (result []*model.DomainContactDetailsHourly, err error) {
	result = make([]*model.DomainContactDetailsHourly, 0)
	for rows.Next() {
		var m = new(model.DomainContactDetailsHourly)
		err = rows.Scan(&m.DomainId,
			&m.DomainName,
			&m.DomainCreatedOn,
			&m.RegistrarExtId,
			&m.RegistrantClientId,
			&m.RegistrantName,
			&m.RegistrantOrg,
			&m.RegistrantEmail)
		if err != nil {
			return
		}
		result = append(result, m)
	}
	err = rows.Err()
	return
}

func (repo *DomainContactDetailsHouryRepository) FindAll() (result []*model.DomainContactDetailsHourly, err error) {
	rows, err := repo.db.Query("SELECT " + repo.FIELDS + " FROM " + repo.TABLE_NAME + " ORDER BY domain_created_on ASC")
	if err != nil {
		return
	}
	defer rows.Close()
	result, err = repo.rowsToResult(rows)
	return
}

func (repo *DomainContactDetailsHouryRepository) Stats() (count int, maxKey string, err error) {
	err = repo.db.QueryRow("SELECT COUNT(domain_name), MAX(domain_id) FROM "+repo.TABLE_NAME).Scan(&count, &maxKey)
	return
}

func (repo *DomainContactDetailsHouryRepository) FindPaginated(numitems int, offsetKey string) (result []*model.DomainContactDetailsHourly, err error) {
	var rows *sql.Rows
	if len(offsetKey) > 0 {
		rows, err = repo.db.Query("SELECT "+repo.FIELDS+" "+"FROM "+repo.TABLE_NAME+" WHERE domain_id > $1 ORDER BY domain_created_on ASC LIMIT $2", offsetKey, numitems)
	} else {
		rows, err = repo.db.Query("SELECT "+repo.FIELDS+" "+"FROM "+repo.TABLE_NAME+" ORDER BY domain_created_on ASC LIMIT $1", numitems)
	}
	if err != nil {
		return
	}
	defer rows.Close()
	result, err = repo.rowsToResult(rows)
	return
}
