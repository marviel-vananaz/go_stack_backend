package repo

import (
	"database/sql"

	. "github.com/go-jet/jet/v2/sqlite"
	"github.com/marviel-vananaz/go-stack-backend/internal/db/model"
	. "github.com/marviel-vananaz/go-stack-backend/internal/db/table"
)

type petRepo struct {
	db *sql.DB
}

func (r *petRepo) Add(name string) *model.Pets {
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
	return &dest
}

func (r *petRepo) Delete(id int) error {
	stmt := Pets.DELETE().WHERE(Pets.ID.EQ(Int(int64(id))))
	_, err := stmt.Exec(r.db)
	if err != nil {
		return err
	}
	return nil
}

func (r *petRepo) GetByID(id int) (*model.Pets, error) {
	stmt := Pets.SELECT(Pets.AllColumns).WHERE(Pets.ID.EQ(Int(int64(id))))
	dest := model.Pets{}
	err := stmt.Query(r.db, &dest)
	if err != nil {
		return nil, err
	}
	return &dest, nil
}

func (r *petRepo) Update(pet *model.Pets) error {
	stmt := Pets.UPDATE(
		Pets.Name,
		Pets.Status,
	).SET(
		pet.Name,
		pet.Status,
	).WHERE(Pets.ID.EQ(Int(int64(*pet.ID))))

	_, err := stmt.Exec(r.db)
	return err
}

func (r *petRepo) List(status *string) ([]*model.Pets, error) {
	stmt := Pets.SELECT(Pets.AllColumns)
	if status != nil {
		stmt = stmt.WHERE(Pets.Status.EQ(String(*status)))
	}

	var dest []*model.Pets
	err := stmt.Query(r.db, &dest)
	if err != nil {
		return nil, err
	}
	return dest, nil
}

func NewPetRepo(db *sql.DB) petRepo {
	return petRepo{
		db: db,
	}
}
