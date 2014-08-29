package registrations

import (
	"database/sql"
	"fmt"
	"github.com/dothiv/afilias-registry-operator-reports/afilias/model"
	"github.com/dothiv/afilias-registry-operator-reports/report"
	"github.com/dothiv/afilias-registry-operator-reports/repository"
	_ "github.com/lib/pq"
	"os"
	"path/filepath"
	"time"
)

var repo *repository.DomainContactDetailsHourlyRepository
var latestImport time.Time

func Import(c *Config) (err error) {
	log(!c.Quiet, "Importing from "+c.ReportsDir)
	log(!c.Quiet, fmt.Sprintf("Importing into %s@%s/%s", c.Database.User, c.Database.Host, c.Database.Name))

	// Open DB
	db, err := sql.Open("postgres", fmt.Sprintf("user=%s dbname=%s sslmode=disable host=%s", c.Database.User, c.Database.Name, c.Database.Host))
	if err != nil {
		return
	}

	// How much can we skip?
	repo = repository.NewDomainContactDetailsHourlyRepository(db)
	latestImportTime, err := repo.GetLatestImportTime()
	if err == sql.ErrNoRows {
		err = nil
		log(!c.Quiet, "Initial import")
	} else {
		log(!c.Quiet, "Importing events after "+latestImportTime)
		latestImport, err = time.Parse("2006-01-02 15:04:05.999999999-07", latestImportTime)
		if err != nil {
			return
		}
	}

	err = filepath.Walk(c.ReportsDir, importFile)

	return
}

func log(b bool, msg string) {
	if b {
		os.Stdout.WriteString(msg + "\n")
	}
}

func importFile(path string, f os.FileInfo, e error) (err error) {
	if f.IsDir() {
		return nil
	}
	var r *model.Report
	r, err = model.NewReportFromName(f.Name())
	if err != nil {
		return nil
	}
	if r.Content != "domain-contact-details" || r.Grouping != "hiv" || r.Frequency != "hourly" {
		return nil
	}

	if &latestImport != nil {
		var reportTime time.Time
		reportTime, err = time.Parse("2006-01-02T15", r.Date)
		if err != nil {
			return
		}
		if latestImport.After(reportTime) {
			return
		}
	}

	reader := report.NewReportReader()
	err = reader.Open(filepath.Dir(path)+"/", r)
	if err != nil {
		return
	}
	defer reader.Close()
	for {
		line, readErr := reader.Next()
		if readErr != nil {
			break
		}
		var domain_created_on time.Time
		domain_created_on, err = time.Parse("2006-01-02 15:04:05.999999999-07", line["domain_created_on"])
		if err != nil {
			return
		}
		if latestImport.After(domain_created_on) || latestImport.Equal(domain_created_on) {
			continue
		}
		dcdh := new(model.DomainContactDetailsHourly)
		dcdh.DomainId = line["domain_id"]
		dcdh.DomainName = line["domain_name"]
		dcdh.DomainCreatedOn = line["domain_created_on"]
		dcdh.RegistrarExtId = line["registrar_ext_id"]
		dcdh.RegistrantClientId = line["registrant_client_id"]
		dcdh.RegistrantName = line["registrant_name"]
		dcdh.RegistrantOrg = line["registrant_org"]
		dcdh.RegistrantEmail = line["registrant_email"]
		_, err = repo.Persist(dcdh)
		if err != nil {
			return
		}
	}

	return
}
