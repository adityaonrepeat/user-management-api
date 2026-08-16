package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	db "github.com/adityaonrepeat/user-management-api/db/sqlc"
	"github.com/adityaonrepeat/user-management-api/internal/models"
)

var ErrNotFound = errors.New("record not found")

type UserRepository struct {
	q *db.Queries
}

func NewUserRepository(q *db.Queries) *UserRepository {
	return &UserRepository{q: q}
}

func (r *UserRepository) Create(ctx context.Context, name string, dob time.Time) (models.User, error) {
	u, err := r.q.CreateUser(ctx, db.CreateUserParams{Name: name, Dob: dob})
	if err != nil {
		return models.User{}, fmt.Errorf("create user: %w", err)
	}
	return toModel(u), nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int32) (models.User, error) {
	u, err := r.q.GetUser(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return models.User{}, ErrNotFound
	}
	if err != nil {
		return models.User{}, fmt.Errorf("get user %d: %w", id, err)
	}
	return toModel(u), nil
}

func (r *UserRepository) List(ctx context.Context) ([]models.User, error) {
	rows, err := r.q.ListUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}

	users := make([]models.User, 0, len(rows))
	for _, u := range rows {
		users = append(users, toModel(u))
	}
	return users, nil
}

func (r *UserRepository) ListPaginated(ctx context.Context, limit, offset int32) ([]models.User, error) {
	rows, err := r.q.ListUsersPaginated(ctx, db.ListUsersPaginatedParams{Limit: limit, Offset: offset})
	if err != nil {
		return nil, fmt.Errorf("list users paginated: %w", err)
	}

	users := make([]models.User, 0, len(rows))
	for _, u := range rows {
		users = append(users, toModel(u))
	}
	return users, nil
}

func (r *UserRepository) Count(ctx context.Context) (int64, error) {
	total, err := r.q.CountUsers(ctx)
	if err != nil {
		return 0, fmt.Errorf("count users: %w", err)
	}
	return total, nil
}

func (r *UserRepository) Update(ctx context.Context, id int32, name string, dob time.Time) (models.User, error) {
	u, err := r.q.UpdateUser(ctx, db.UpdateUserParams{ID: id, Name: name, Dob: dob})
	if errors.Is(err, pgx.ErrNoRows) {
		return models.User{}, ErrNotFound
	}
	if err != nil {
		return models.User{}, fmt.Errorf("update user %d: %w", id, err)
	}
	return toModel(u), nil
}

func (r *UserRepository) Delete(ctx context.Context, id int32) error {
	affected, err := r.q.DeleteUser(ctx, id)
	if err != nil {
		return fmt.Errorf("delete user %d: %w", id, err)
	}
	if affected == 0 {
		return ErrNotFound
	}
	return nil
}

func toModel(u db.User) models.User {
	return models.User{ID: u.ID, Name: u.Name, DOB: u.Dob}
}
