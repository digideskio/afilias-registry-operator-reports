package repository

import (
	"database/sql"
	"github.com/dothiv/afilias-registry-operator-reports/afilias/model"
	_ "github.com/lib/pq"
)

type TransactionRepositoryInterface interface {
	GetLatestImportTime() (transaction_date string, err error)
	Persist(m *model.Transaction) (result sql.Result, err error)
	FindAll() (result []*model.Transaction, err error)
	Stats() (count int, maxKey string, err error)
	FindPaginated(numitems int, offsetKey string) (result []*model.Transaction, err error)
}

type TransactionRepository struct {
	TransactionRepositoryInterface
	db         *sql.DB
	TABLE_NAME string
	FIELDS     string
}

func NewTransactionRepository(db *sql.DB) (repo *TransactionRepository) {
	repo = new(TransactionRepository)
	repo.db = db
	repo.TABLE_NAME = "transaction"
	repo.FIELDS = "tld, registrar_ext_id, registrar_name, server_transaction_id, command, object_type, object_name, transaction_date"
	return
}

func (repo *TransactionRepository) GetLatestImportTime() (transaction_date string, err error) {
	err = repo.db.QueryRow("SELECT transaction_date FROM " + repo.TABLE_NAME + " ORDER BY server_transaction_id DESC LIMIT 1").Scan(&transaction_date)
	return
}

func (repo *TransactionRepository) Persist(m *model.Transaction) (result sql.Result, err error) {
	result, err = repo.db.Exec("INSERT INTO "+repo.TABLE_NAME+" "+
		"("+repo.FIELDS+") "+
		"VALUES($1, $2, $3, $4, $5, $6, $7, $8)",
		m.TLD,
		m.Registrar_Ext_ID,
		m.Registrar_Name,
		m.Server_TrID,
		m.Command,
		m.Object_Type,
		m.Object_Name,
		m.Transaction_Date)
	return
}

func (repo *TransactionRepository) rowsToResult(rows *sql.Rows) (result []*model.Transaction, err error) {
	result = make([]*model.Transaction, 0)
	for rows.Next() {
		var m = new(model.Transaction)
		err = rows.Scan(
			&m.TLD,
			&m.Registrar_Ext_ID,
			&m.Registrar_Name,
			&m.Server_TrID,
			&m.Command,
			&m.Object_Type,
			&m.Object_Name,
			&m.Transaction_Date)
		if err != nil {
			return
		}
		result = append(result, m)
	}
	err = rows.Err()
	return
}

func (repo *TransactionRepository) FindAll() (result []*model.Transaction, err error) {
	rows, err := repo.db.Query("SELECT " + repo.FIELDS + " FROM " + repo.TABLE_NAME + " ORDER BY server_transaction_id ASC")
	if err != nil {
		return
	}
	defer rows.Close()
	result, err = repo.rowsToResult(rows)
	return
}

func (repo *TransactionRepository) Stats() (count int, maxKey string, err error) {
	err = repo.db.QueryRow("SELECT COUNT(server_transaction_id), MAX(server_transaction_id) FROM "+repo.TABLE_NAME).Scan(&count, &maxKey)
	return
}

func (repo *TransactionRepository) FindPaginated(numitems int, offsetKey string) (result []*model.Transaction, err error) {
	var rows *sql.Rows
	if len(offsetKey) > 0 {
		rows, err = repo.db.Query("SELECT "+repo.FIELDS+" "+"FROM "+repo.TABLE_NAME+" WHERE server_transaction_id > $1 ORDER BY server_transaction_id ASC LIMIT $2", offsetKey, numitems)
	} else {
		rows, err = repo.db.Query("SELECT "+repo.FIELDS+" "+"FROM "+repo.TABLE_NAME+" ORDER BY server_transaction_id ASC LIMIT $1", numitems)
	}
	if err != nil {
		return
	}
	defer rows.Close()
	result, err = repo.rowsToResult(rows)
	return
}
