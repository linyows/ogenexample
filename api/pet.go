package api

import (
	"context"
	"net/http"

	"github.com/linyows/ogenexample/db/dbgen"
	"github.com/linyows/ogenexample/oas/oasgen"
)

func (h *handler) AddPet(ctx context.Context, req *oasgen.Pet) (*oasgen.Pet, error) {
	tracer := h.tp.Tracer("db-trace")
	ctx, span := tracer.Start(ctx, "sqlc.CreatePet")
	defer span.End()

	st := oasgen.PetStatus("pending")
	s := req.GetStatus()
	if s.Set {
		st = s.Value
	}
	res, err := h.q.CreatePet(ctx, dbgen.CreatePetParams{
		Name:   req.GetName(),
		Status: dbgen.PetsStatus(st),
	})
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &oasgen.Pet{
		ID:     oasgen.NewOptInt64(id),
		Name:   req.GetName(),
		Status: oasgen.NewOptPetStatus(oasgen.PetStatus(st)),
	}, nil
}

func (h *handler) DeletePet(ctx context.Context, params oasgen.DeletePetParams) error {
	tracer := h.tp.Tracer("db-trace")
	ctx, span := tracer.Start(ctx, "sqlc.DeletePet")
	defer span.End()

	if err := h.q.DeletePet(ctx, params.PetId); err != nil {
		return err
	}
	return nil
}

func (h *handler) GetPetById(ctx context.Context, params oasgen.GetPetByIdParams) (oasgen.GetPetByIdRes, error) {
	tracer := h.tp.Tracer("db-trace")
	ctx, span := tracer.Start(ctx, "sqlc.GetPetByID")
	defer span.End()

	pet, err := h.q.GetPetByID(ctx, params.PetId)
	if err != nil {
		return &oasgen.GetPetByIdNotFound{}, nil
	}
	return &oasgen.Pet{
		ID:     oasgen.NewOptInt64(pet.ID),
		Name:   pet.Name,
		Status: oasgen.NewOptPetStatus(oasgen.PetStatus(pet.Status)),
	}, nil
}

func (h *handler) UpdatePet(ctx context.Context, params oasgen.UpdatePetParams) error {
	tracer := h.tp.Tracer("db-trace")
	ctx, span := tracer.Start(ctx, "sqlc.UpdatePet")
	defer span.End()

	return h.q.UpdatePet(ctx, dbgen.UpdatePetParams{
		Name:   params.Name.Value,
		Status: dbgen.PetsStatus(params.Status.Value),
		ID:     params.PetId,
	})
}

func (h *handler) NewError(ctx context.Context, err error) (r *oasgen.ErrorStatusCode) {
	return &oasgen.ErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: oasgen.Error{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		},
	}
}
