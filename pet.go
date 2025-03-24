package main

import (
	"context"
	"errors"
	"net/http"

	"github.com/marviel-vananaz/go-stack-backend/internal/db/model"
	"github.com/marviel-vananaz/go-stack-backend/internal/oas"
	"github.com/marviel-vananaz/go-stack-backend/internal/repo"
)

type petRepo interface {
	// Add creates a new pet with the given name
	Add(name string) (*model.Pets, error)

	// Delete removes a pet by its ID
	Delete(id int) error

	// GetByID retrieves a pet by its ID
	GetByID(id int) (*model.Pets, error)

	// Update modifies an existing pet
	Update(pet *model.Pets) error

	// List retrieves all pets, optionally filtered by status
	List(status *string) ([]*model.Pets, error)
}

type petsService struct {
	repo petRepo
}

func (p *petsService) AddPet(ctx context.Context, req *oas.Pet) (oas.AddPetRes, error) {
	if req.Name == "" {
		res := oas.AddPetBadRequest{
			Code:    http.StatusBadRequest,
			Message: "pet name is required",
		}
		return &res, nil
	}

	pet, err := p.repo.Add(req.Name)
	if err != nil {
		res := oas.AddPetInternalServerError{
			Code:    http.StatusInternalServerError,
			Message: "failed to create pet",
		}
		return &res, nil
	}

	return &oas.Pet{
		ID:     oas.NewOptInt64(int64(*pet.ID)),
		Name:   pet.Name,
		Status: oas.NewOptPetStatus(oas.PetStatus(*pet.Status)),
	}, nil
}

func (p *petsService) DeletePet(ctx context.Context, params oas.DeletePetParams) (oas.DeletePetRes, error) {
	err := p.repo.Delete(int(params.PetId))
	if err != nil {
		if errors.Is(err, repo.ErrPetNotFound) {
			res := oas.DeletePetNotFound{
				Code:    http.StatusNotFound,
				Message: "pet not found",
			}
			return &res, nil
		}
		res := oas.DeletePetInternalServerError{
			Code:    http.StatusInternalServerError,
			Message: "failed to delete pet",
		}
		return &res, nil
	}
	return &oas.DeletePetOK{}, nil
}

func (p *petsService) GetPetById(ctx context.Context, params oas.GetPetByIdParams) (oas.GetPetByIdRes, error) {
	pet, err := p.repo.GetByID(int(params.PetId))
	if err != nil {
		if errors.Is(err, repo.ErrPetNotFound) {
			res := oas.GetPetByIdNotFound{
				Code:    http.StatusNotFound,
				Message: "pet not found",
			}
			return &res, nil
		}
		res := oas.GetPetByIdInternalServerError{
			Code:    http.StatusInternalServerError,
			Message: "failed to get pet",
		}
		return &res, nil
	}

	return &oas.Pet{
		ID:     oas.NewOptInt64(int64(*pet.ID)),
		Name:   pet.Name,
		Status: oas.NewOptPetStatus(oas.PetStatus(*pet.Status)),
	}, nil
}

func (p *petsService) UpdatePet(ctx context.Context, params oas.UpdatePetParams) (oas.UpdatePetRes, error) {
	if !params.Name.Set || params.Name.Value == "" {
		res := oas.UpdatePetBadRequest{
			Code:    http.StatusBadRequest,
			Message: "pet name is required",
		}
		return &res, nil
	}

	id := int32(params.PetId)
	err := p.repo.Update(&model.Pets{
		ID:     &id,
		Name:   params.Name.Value,
		Status: (*string)(&params.Status.Value),
	})
	if err != nil {
		if errors.Is(err, repo.ErrPetNotFound) {
			res := oas.UpdatePetNotFound{
				Code:    http.StatusNotFound,
				Message: "pet not found",
			}
			return &res, nil
		}
		res := oas.UpdatePetInternalServerError{
			Code:    http.StatusInternalServerError,
			Message: "failed to update pet",
		}
		return &res, nil
	}
	return nil, nil
}

func (p *petsService) ListPets(ctx context.Context) (oas.ListPetsRes, error) {
	pets, err := p.repo.List(nil)
	if err != nil {
		res := oas.ListPetsInternalServerError{
			Code:    http.StatusInternalServerError,
			Message: "failed to list pets",
		}
		return &res, nil
	}

	result := make([]oas.Pet, len(pets))
	for i, pet := range pets {
		result[i] = oas.Pet{
			ID:     oas.NewOptInt64(int64(*pet.ID)),
			Name:   pet.Name,
			Status: oas.NewOptPetStatus(oas.PetStatus(*pet.Status)),
		}
	}

	res := oas.ListPetsOKApplicationJSON(result)
	return &res, nil
}
