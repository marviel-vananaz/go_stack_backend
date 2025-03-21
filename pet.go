package main

import (
	"context"

	"github.com/marviel-vananaz/go-stack-backend/internal/db/model"
	"github.com/marviel-vananaz/go-stack-backend/internal/oas"
)

type petRepo interface {
	Add(name string) model.Pets
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
	return nil
}

func (p *petsService) GetPetById(ctx context.Context, params oas.GetPetByIdParams) (oas.GetPetByIdRes, error) {
	return nil, nil
}

func (p *petsService) UpdatePet(ctx context.Context, params oas.UpdatePetParams) error {

	return nil
}

func (p *petsService) ListPets(ctx context.Context) ([]oas.Pet, error) {
	return nil, nil
}
