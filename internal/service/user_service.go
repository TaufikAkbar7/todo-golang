package service

import (
	"context"
	"database/sql"
	"errors"
	"golang-todo/internal/entity"
	"golang-todo/internal/helper"
	"golang-todo/internal/model"
	"golang-todo/internal/repository"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	DB        *sqlx.DB
	Repo      *repository.UserRepository
	Log       *logrus.Logger
	Validator *validator.Validate
}

func NewUserService(db *sqlx.DB, repo *repository.UserRepository, log *logrus.Logger, validator *validator.Validate) *UserService {
	return &UserService{
		DB:        db,
		Repo:      repo,
		Log:       log,
		Validator: validator,
	}
}

func (c *UserService) Create(ctx context.Context, req *model.UserCreateRequest) error {
	newID, _ := uuid.NewV7()
	tx, _ := c.DB.BeginTxx(ctx, nil)
	defer tx.Rollback()

	if err := c.Validator.Struct(req); err != nil {
		c.Log.Errorf("Invalid request body  : %+v", err)
		return fiber.ErrBadRequest
	}

	// find user
	user, err := c.Repo.FindByUsername(ctx, tx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows || err == errors.New("not Found") {
			// generate pw
			password, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
			if err != nil {
				c.Log.Errorf("Failed hash password %v", err)
				return fiber.ErrInternalServerError
			}

			// store on db
			dateNow := helper.GetDateNow()
			newUser := &entity.User{
				ID:        newID,
				Username:  req.Username,
				Password:  string(password),
				CreatedAt: dateNow,
				UpdatedAt: dateNow,
			}
			if err := c.Repo.Create(ctx, tx, newUser); err != nil {
				c.Log.Errorf("Failed insert data to db %v", err)
				return fiber.ErrInternalServerError
			}

			if err := tx.Commit(); err != nil {
				c.Log.Errorf("Failed commit transaction : %v", err)
				return fiber.ErrInternalServerError
			}

			return nil
		}

		c.Log.Errorf("Failed search user by id %v", err)
		return fiber.ErrInternalServerError
	}
	// user already exist
	if user != nil {
		return fiber.NewError(400, "Already registered user")
	}

	c.Log.Info("Success created user")
	return nil
}

func (c *UserService) Login(ctx context.Context, req *model.LoginUserRequest) (*model.UserReponseLogin, error) {
	tx, _ := c.DB.BeginTxx(ctx, nil)
	defer tx.Rollback()

	if err := c.Validator.Struct(req); err != nil {
		c.Log.Errorf("Invalid request body  : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	user, err := c.Repo.FindByUsername(ctx, tx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fiber.ErrUnauthorized
		}
		c.Log.Errorf("Failed search user by id %v", err)
		return nil, fiber.ErrUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.Log.Errorf("Failed to compare user password with bcrype hash : %+v", err)
		return nil, fiber.ErrUnauthorized
	}

	expirationTime := time.Now().Add(1 * time.Hour) // set expiration time for 1 hour
	userClaims := model.UserCustomClaims{
		Id:       user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime), // Set the expiration time
			IssuedAt:  jwt.NewNumericDate(time.Now()),     // Optional: Set issued at time
			Issuer:    "todo-golang",                      // Optional: Set issuer
		},
	}

	token := helper.GenerateJWTToken(userClaims, c.Log)
	entity := &entity.User{
		ID:        user.ID,
		Password:  user.Password,
		Username:  user.Username,
		UpdatedAt: helper.GetDateNow(),
	}

	if err := c.Repo.Update(ctx, tx, entity); err != nil {
		c.Log.Errorf("Failed save user : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit(); err != nil {
		c.Log.Errorf("Failed commit transaction : %v", err)
		return nil, fiber.ErrInternalServerError
	}

	response := &model.UserReponseLogin{
		Token: token,
	}

	c.Log.Info("Success login user")
	return response, nil
}

func (c *UserService) ChangePassword(ctx context.Context, req *model.UserChangePasswordRequest, id uuid.UUID) error {
	tx, err := c.DB.BeginTxx(ctx, nil)
	if err != nil {
		c.Log.Errorf("Failed start transaction db %v", err)
		return err
	}

	user, err := c.Repo.FindByID(ctx, tx, id)
	if err != nil {
		c.Log.Errorf("Failed to get data user by id %v", err)
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.Log.Errorf("Failed compare hash password %v", err)
		return err
	}

	newPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Log.Errorf("Failed hash new password %v", err)
		return err
	}

	payload := model.UserChangePassword{
		Password: string(newPassword),
		UpdateAt: helper.GetDateNow(),
	}
	if err := c.Repo.ChangePassword(ctx, tx, payload, id); err != nil {
		c.Log.Errorf("Failed to change password %v", err)
		return err
	}

	c.Log.Info("Success change password")
	return nil
}

func (c *UserService) Verify(ctx context.Context, req *model.VerifyUserRequest) (*model.Auth, error) {
	tx, _ := c.DB.BeginTxx(ctx, nil)
	defer tx.Rollback()

	claim, err := helper.ParseJWTToken(c.Log, req.Token, &model.UserCustomClaims{})
	if claim != nil {
		// check if user trusted
		user, err := c.Repo.FindByID(ctx, tx, claim.Id)
		if err != nil {
			c.Log.Errorf("Failed find user by id:%v", err)
			return nil, fiber.ErrUnauthorized
		}
		auth := new(model.Auth)
		auth.ID = user.ID.String()
		auth.Username = user.Username
		return auth, nil
	}

	c.Log.Info("Success verify user by token")
	return nil, err
}
