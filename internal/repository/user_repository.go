package repository

import (
	"context"
	"golang-todo/internal/entity"

	"golang-todo/internal/model"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type UserRepository struct {
	DB  *sqlx.DB
	Log *logrus.Logger
}

func NewUserRepository(db *sqlx.DB, log *logrus.Logger) *UserRepository {
	return &UserRepository{
		DB:  db,
		Log: log,
	}
}

func (userRepo *UserRepository) Create(ctx context.Context, tx *sqlx.Tx, entity *entity.User) error {
	query := "INSERT INTO users (id, password, username) VALUES ($1, $2, $3)"
	_, err := tx.ExecContext(ctx, query, entity.ID, entity.Password, entity.Username)
	if err != nil {
		return err
	}
	return nil
}

func (userRepo *UserRepository) ChangePassword(ctx context.Context, tx *sqlx.Tx, user model.UserChangePassword, id uuid.UUID) error {
	query := "UPDATE users SET password = $1, updated_at = $2 WHERE id = $3"

	_, err := tx.ExecContext(ctx, query, user.Password, user.UpdateAt, id)
	return err
}

func (userRepo *UserRepository) FindByID(ctx context.Context, tx *sqlx.Tx, id uuid.UUID) (*entity.User, error) {
	query := "SELECT id, username FROM users WHERE id = $1"
	user := new(entity.User)

	err := tx.GetContext(ctx, user, query, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (userRepo *UserRepository) FindByUsername(ctx context.Context, tx *sqlx.Tx, name string) (*model.UserResultFindUsername, error) {
	query := "SELECT id, username, password FROM users WHERE username = $1"
	user := new(model.UserResultFindUsername)
	err := tx.GetContext(ctx, user, query, name)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (userRepo *UserRepository) Update(ctx context.Context, tx *sqlx.Tx, entity *entity.User) error {
	query := "UPDATE users SET password = $1, username = $2, updated_at = $3 WHERE id = $4"

	_, err := tx.ExecContext(ctx, query, entity.Password, entity.Username, entity.UpdatedAt, entity.ID)
	return err
}
