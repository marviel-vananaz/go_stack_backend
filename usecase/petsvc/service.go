package petsvc

import (
	"context"
	"errors"
	"net/http"

	"github.com/marviel-vananaz/go-stack-backend/.gen/api"
	"github.com/marviel-vananaz/go-stack-backend/.gen/db/model"
	"github.com/marviel-vananaz/go-stack-backend/infra/sqlite"
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

type petService struct {
	repo petRepo
}

func (p *petService) AddPet(ctx context.Context, req *api.Pet) (api.AddPetRes, error) {
	if req.Name == "" {
		res := api.AddPetBadRequest{
			Code:    http.StatusBadRequest,
			Message: "pet name is required",
		}
		return &res, nil
	}

	pet, err := p.repo.Add(req.Name)
	if err != nil {
		res := api.AddPetInternalServerError{
			Code:    http.StatusInternalServerError,
			Message: "failed to create pet",
		}
		return &res, nil
	}

	return &api.Pet{
		ID:     api.NewOptInt64(int64(*pet.ID)),
		Name:   pet.Name,
		Status: api.NewOptPetStatus(api.PetStatus(*pet.Status)),
	}, nil
}

func (p *petService) DeletePet(ctx context.Context, params api.DeletePetParams) (api.DeletePetRes, error) {
	err := p.repo.Delete(int(params.PetId))
	if err != nil {
		if errors.Is(err, sqlite.ErrPetNotFound) {
			res := api.DeletePetNotFound{
				Code:    http.StatusNotFound,
				Message: "pet not found",
			}
			return &res, nil
		}
		res := api.DeletePetInternalServerError{
			Code:    http.StatusInternalServerError,
			Message: "failed to delete pet",
		}
		return &res, nil
	}
	return &api.DeletePetOK{}, nil
}

func (p *petService) GetPetById(ctx context.Context, params api.GetPetByIdParams) (api.GetPetByIdRes, error) {
	pet, err := p.repo.GetByID(int(params.PetId))
	if err != nil {
		if errors.Is(err, sqlite.ErrPetNotFound) {
			res := api.GetPetByIdNotFound{
				Code:    http.StatusNotFound,
				Message: "pet not found",
			}
			return &res, nil
		}
		res := api.GetPetByIdInternalServerError{
			Code:    http.StatusInternalServerError,
			Message: "failed to get pet",
		}
		return &res, nil
	}

	return &api.Pet{
		ID:     api.NewOptInt64(int64(*pet.ID)),
		Name:   pet.Name,
		Status: api.NewOptPetStatus(api.PetStatus(*pet.Status)),
	}, nil
}

func (p *petService) UpdatePet(ctx context.Context, params api.UpdatePetParams) (api.UpdatePetRes, error) {
	if !params.Name.Set || params.Name.Value == "" {
		res := api.UpdatePetBadRequest{
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
		if errors.Is(err, sqlite.ErrPetNotFound) {
			res := api.UpdatePetNotFound{
				Code:    http.StatusNotFound,
				Message: "pet not found",
			}
			return &res, nil
		}
		res := api.UpdatePetInternalServerError{
			Code:    http.StatusInternalServerError,
			Message: "failed to update pet",
		}
		return &res, nil
	}
	return nil, nil
}

func (p *petService) ListPets(ctx context.Context) (api.ListPetsRes, error) {
	pets, err := p.repo.List(nil)
	if err != nil {
		res := api.ListPetsInternalServerError{
			Code:    http.StatusInternalServerError,
			Message: "failed to list pets",
		}
		return &res, nil
	}

	result := make([]api.Pet, len(pets))
	for i, pet := range pets {
		result[i] = api.Pet{
			ID:     api.NewOptInt64(int64(*pet.ID)),
			Name:   pet.Name,
			Status: api.NewOptPetStatus(api.PetStatus(*pet.Status)),
		}
	}

	res := api.ListPetsOKApplicationJSON(result)
	return &res, nil
}

func NewService(repo petRepo) *petService {
	return &petService{
		repo: repo,
	}
}
