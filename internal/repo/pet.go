package repo

import (
	"database/sql"
	"errors"
	"fmt"

	. "github.com/go-jet/jet/v2/sqlite"
	"github.com/marviel-vananaz/go-stack-backend/internal/db/model"
	. "github.com/marviel-vananaz/go-stack-backend/internal/db/table"
)

var (
	ErrPetNotFound = errors.New("pet not found")
)

type petRepo struct {
	db *sql.DB
}

func (r *petRepo) Add(name string) (*model.Pets, error) {
	status := "available"
	stmt := Pets.INSERT(Pets.AllColumns).MODEL(&model.Pets{
		Name:   name,
		Status: &status,
	}).RETURNING(Pets.AllColumns)
	dest := model.Pets{}
	err := stmt.Query(r.db, &dest)
	if err != nil {
		return nil, fmt.Errorf("failed to add pet: %w", err)
	}
	return &dest, nil
}

func (r *petRepo) Delete(id int) error {
	stmt := Pets.DELETE().WHERE(Pets.ID.EQ(Int(int64(id))))
	result, err := stmt.Exec(r.db)
	if err != nil {
		return fmt.Errorf("failed to delete pet: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return ErrPetNotFound
	}
	return nil
}

func (r *petRepo) GetByID(id int) (*model.Pets, error) {
	stmt := Pets.SELECT(Pets.AllColumns).WHERE(Pets.ID.EQ(Int(int64(id))))
	dest := model.Pets{}
	err := stmt.Query(r.db, &dest)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrPetNotFound
		}
		return nil, fmt.Errorf("failed to get pet: %w", err)
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

	result, err := stmt.Exec(r.db)
	if err != nil {
		return fmt.Errorf("failed to update pet: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return ErrPetNotFound
	}
	return nil
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
