package transactions

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

var repo repository.TransactionRepositoryInterface
var latestImport time.Time

func Import(c *Config) (err error) {
	log(!c.Quiet, "Importing from "+c.ReportsDir)
	log(!c.Quiet, fmt.Sprintf("Importing into %s@%s/%s", c.Database.User, c.Database.Host, c.Database.Name))

	// Open DB
	db, err := sql.Open("postgres", c.Database.DSN())
	if err != nil {
		return
	}

	// How much can we skip?
	repo = repository.NewTransactionRepository(db)
	latestImportTime, err := repo.GetLatestImportTime()
	if err == sql.ErrNoRows {
		err = nil
		log(!c.Quiet, "Initial import")
	} else if err != nil {
		return
	} else {
		log(!c.Quiet, "Importing events after "+latestImportTime)
		latestImport, err = time.Parse("2006-01-02 15:04:05", latestImportTime)
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
	if r.Content != "transactions" || r.Grouping != "ALL" || r.Frequency != "daily" {
		return nil
	}

	if &latestImport != nil {
		var reportTime time.Time
		reportTime, err = time.Parse("2006-01-02", r.Date)
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
		var transaction_date time.Time
		transaction_date, err = time.Parse("2006-01-02 15:04:05", line["Transaction_Date"])
		if err != nil {
			return
		}
		if latestImport.After(transaction_date) || latestImport.Equal(transaction_date) {
			continue
		}
		transaction := new(model.Transaction)
		transaction.TLD = line["TLD"]
		transaction.Registrar_Ext_ID = line["Registrar_Ext_ID"]
		transaction.Registrar_Name = line["Registrar_Name"]
		transaction.Server_TrID = line["Server_TrID"]
		transaction.Command = line["Command"]
		transaction.Object_Type = line["Object_Type"]
		transaction.Object_Name = line["Object_Name"]
		transaction.Transaction_Date = line["Transaction_Date"]
		_, err = repo.Persist(transaction)
		if err != nil {
			return
		}
	}

	return
}
