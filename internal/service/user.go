package service

import (
	"context"
	"errors"
	"time"

	"github.com/adityaonrepeat/user-management-api/internal/models"
	"github.com/adityaonrepeat/user-management-api/internal/repository"
)

var ErrInvalidDOB = errors.New("dob must be a valid date in YYYY-MM-DD format")

var ErrDOBInFuture = errors.New("dob cannot be in the future")

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Create(ctx context.Context, req models.CreateUserRequest) (models.UserResponse, error) {
	dob, err := parseDOB(req.DOB)
	if err != nil {
		return models.UserResponse{}, err
	}

	u, err := s.repo.Create(ctx, req.Name, dob)
	if err != nil {
		return models.UserResponse{}, err
	}
	return toResponse(u), nil
}

func (s *UserService) GetByID(ctx context.Context, id int32) (models.UserWithAgeResponse, error) {
	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return models.UserWithAgeResponse{}, err
	}
	return toResponseWithAge(u, time.Now().UTC()), nil
}

func (s *UserService) List(ctx context.Context, params models.ListParams) ([]models.UserWithAgeResponse, int64, error) {
	var (
		users []models.User
		total int64
		err   error
	)

	if params.Paginated {
		if users, err = s.repo.ListPaginated(ctx, params.Limit, params.Offset); err != nil {
			return nil, 0, err
		}
		if total, err = s.repo.Count(ctx); err != nil {
			return nil, 0, err
		}
	} else {
		if users, err = s.repo.List(ctx); err != nil {
			return nil, 0, err
		}
		total = int64(len(users))
	}

	now := time.Now().UTC()

	out := make([]models.UserWithAgeResponse, 0, len(users))
	for _, u := range users {
		out = append(out, toResponseWithAge(u, now))
	}
	return out, total, nil
}

func (s *UserService) Update(ctx context.Context, id int32, req models.UpdateUserRequest) (models.UserResponse, error) {
	dob, err := parseDOB(req.DOB)
	if err != nil {
		return models.UserResponse{}, err
	}

	u, err := s.repo.Update(ctx, id, req.Name, dob)
	if err != nil {
		return models.UserResponse{}, err
	}
	return toResponse(u), nil
}

func (s *UserService) Delete(ctx context.Context, id int32) error {
	return s.repo.Delete(ctx, id)
}

func parseDOB(s string) (time.Time, error) {
	dob, err := time.Parse(models.DateLayout, s)
	if err != nil {
		return time.Time{}, ErrInvalidDOB
	}
	if dob.After(time.Now().UTC()) {
		return time.Time{}, ErrDOBInFuture
	}
	return dob, nil
}

func toResponse(u models.User) models.UserResponse {
	return models.UserResponse{
		ID:   u.ID,
		Name: u.Name,
		DOB:  u.DOB.Format(models.DateLayout),
	}
}

func toResponseWithAge(u models.User, now time.Time) models.UserWithAgeResponse {
	return models.UserWithAgeResponse{
		ID:   u.ID,
		Name: u.Name,
		DOB:  u.DOB.Format(models.DateLayout),
		Age:  CalculateAge(u.DOB, now),
	}
}
