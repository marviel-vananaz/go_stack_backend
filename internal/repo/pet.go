package repo

import (
	"database/sql"

	"github.com/marviel-vananaz/go-stack-backend/internal/db/model"
	. "github.com/marviel-vananaz/go-stack-backend/internal/db/table"
)

type petRepo struct {
	db *sql.DB
}

func (r *petRepo) Add(name string) model.Pets {
	status := "available"
	stmt := Pets.INSERT(Pets.AllColumns).MODEL(&model.Pets{
		Name:   name,
		Status: &status,
	}).RETURNING(Pets.AllColumns)
	dest := model.Pets{}
	err := stmt.Query(r.db, &dest)
	if err != nil {
		stmtstr, _ := stmt.Sql()
		panic("err query: " + stmtstr)
	}
	return dest
}

func NewPetRepo(db *sql.DB) petRepo {
	return petRepo{
		db: db,
	}
}
