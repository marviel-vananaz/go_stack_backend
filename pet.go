package main

import (
	"context"
	"sync"

	"github.com/marviel-vananaz/go-stack-backend/internal/oas"
)

type petsService struct {
	pets map[int64]oas.Pet
	id   int64
	mux  sync.Mutex
}

func (p *petsService) AddPet(ctx context.Context, req *oas.Pet) (*oas.Pet, error) {
	p.mux.Lock()
	defer p.mux.Unlock()

	p.pets[p.id] = *req
	p.id++
	return req, nil
}

func (p *petsService) DeletePet(ctx context.Context, params oas.DeletePetParams) error {
	p.mux.Lock()
	defer p.mux.Unlock()

	delete(p.pets, params.PetId)
	return nil
}

func (p *petsService) GetPetById(ctx context.Context, params oas.GetPetByIdParams) (oas.GetPetByIdRes, error) {
	p.mux.Lock()
	defer p.mux.Unlock()

	pet, ok := p.pets[params.PetId]
	if !ok {
		// Return Not Found.
		return &oas.GetPetByIdNotFound{}, nil
	}
	return &pet, nil
}

func (p *petsService) UpdatePet(ctx context.Context, params oas.UpdatePetParams) error {
	p.mux.Lock()
	defer p.mux.Unlock()

	pet := p.pets[params.PetId]
	pet.Status = params.Status
	if val, ok := params.Name.Get(); ok {
		pet.Name = val
	}
	p.pets[params.PetId] = pet

	return nil
}

func (p *petsService) ListPets(ctx context.Context) ([]oas.Pet, error) {
	p.mux.Lock()
	defer p.mux.Unlock()

	pets := make([]oas.Pet, 0, len(p.pets))
	for _, pet := range p.pets {
		pets = append(pets, pet)
	}
	return pets, nil
}
