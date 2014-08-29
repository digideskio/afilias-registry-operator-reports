package repository

import (
	"database/sql"
	"github.com/dothiv/afilias-registry-operator-reports/afilias/model"
)

type Repository struct {
	db *sql.DB
}

func (r *Repository) Persist(m *model.Model) {}
