package main

import (
	"context"

	"github.com/marviel-vananaz/go-stack-backend/internal/db/model"
	"github.com/marviel-vananaz/go-stack-backend/internal/oas"
)

type petRepo interface {
	// Add creates a new pet with the given name
	Add(name string) *model.Pets

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

func (p *petsService) AddPet(ctx context.Context, req *oas.Pet) (*oas.Pet, error) {
	pet := p.repo.Add(req.Name)
	return &oas.Pet{
		ID:     oas.NewOptInt64(int64(*pet.ID)),
		Name:   pet.Name,
		Status: oas.NewOptPetStatus(oas.PetStatus(*pet.Status)),
	}, nil
}

func (p *petsService) DeletePet(ctx context.Context, params oas.DeletePetParams) error {
	return p.repo.Delete(int(params.PetId))
}

func (p *petsService) GetPetById(ctx context.Context, params oas.GetPetByIdParams) (oas.GetPetByIdRes, error) {
	pet, err := p.repo.GetByID(int(params.PetId))
	if err != nil {
		return nil, err
	}

	return &oas.Pet{
		ID:     oas.NewOptInt64(int64(*pet.ID)),
		Name:   pet.Name,
		Status: oas.NewOptPetStatus(oas.PetStatus(*pet.Status)),
	}, nil
}

func (p *petsService) UpdatePet(ctx context.Context, params oas.UpdatePetParams) error {
	id := int32(params.PetId)
	return p.repo.Update(&model.Pets{
		ID:     &id,
		Name:   params.Name.Value,
		Status: (*string)(&params.Status.Value),
	})
}

func (p *petsService) ListPets(ctx context.Context) ([]oas.Pet, error) {
	pets, err := p.repo.List(nil) // Passing nil to get all pets regardless of status
	if err != nil {
		return nil, err
	}

	result := make([]oas.Pet, len(pets))
	for i, pet := range pets {
		result[i] = oas.Pet{
			ID:     oas.NewOptInt64(int64(*pet.ID)),
			Name:   pet.Name,
			Status: oas.NewOptPetStatus(oas.PetStatus(*pet.Status)),
		}
	}
	return result, nil
}
