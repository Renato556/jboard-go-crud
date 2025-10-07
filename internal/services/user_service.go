package services

import (
	"context"
	"errors"
	"jboard-go-crud/internal/models"
	"jboard-go-crud/internal/models/enums"
	"jboard-go-crud/internal/repositories"
	"log"
	"strings"
)

type UserService interface {
	CreateUser(ctx context.Context, username, password string, role enums.RoleEnum) error
	GetUserByID(ctx context.Context, id string) (models.User, error)
	GetUserByUsername(ctx context.Context, username string) (models.User, error)
	UpdateUser(ctx context.Context, username string, password string, role enums.RoleEnum) (models.User, error)
	DeleteUser(ctx context.Context, username string) error
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	log.Printf("Creating new UserService")
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) CreateUser(ctx context.Context, username, password string, role enums.RoleEnum) error {
	log.Printf("Service CreateUser called for username: %s", username)

	if strings.TrimSpace(username) == "" {
		return errors.New("username cannot be empty")
	}
	if strings.TrimSpace(password) == "" {
		return errors.New("password cannot be empty")
	}
	if !role.IsValid() {
		return errors.New("invalid role")
	}

	_, exists, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		log.Printf("Error checking existing username: %v", err)
		return err
	}
	if exists {
		log.Printf("Username already exists: %s", username)
		return errors.New("username already exists")
	}

	user := models.User{
		Username: username,
		Password: password,
		Role:     role,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		log.Printf("Repository error in CreateUser: %v", err)
		return err
	}

	log.Printf("Successfully created user: %s", username)
	return nil
}

func (s *userService) GetUserByID(ctx context.Context, id string) (models.User, error) {
	log.Printf("Service GetUserByID called for ID: %s", id)

	if strings.TrimSpace(id) == "" {
		return models.User{}, errors.New("user ID cannot be empty")
	}

	user, found, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		log.Printf("Repository error in GetUserByID: %v", err)
		return models.User{}, err
	}

	if !found {
		log.Printf("User not found with ID: %s", id)
		return models.User{}, errors.New("user not found")
	}

	log.Printf("Successfully retrieved user ID: %s", id)
	return user, nil
}

func (s *userService) GetUserByUsername(ctx context.Context, username string) (models.User, error) {
	log.Printf("Service GetUserByUsername called for username: %s", username)

	if strings.TrimSpace(username) == "" {
		return models.User{}, errors.New("username cannot be empty")
	}

	user, found, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		log.Printf("Repository error in GetUserByUsername: %v", err)
		return models.User{}, err
	}

	if !found {
		log.Printf("User not found with username: %s", username)
		return models.User{}, errors.New("user not found")
	}

	log.Printf("Successfully retrieved user by username: %s", username)
	return user, nil
}

func (s *userService) UpdateUser(ctx context.Context, username string, password string, role enums.RoleEnum) (models.User, error) {
	log.Printf("Service UpdateUser called for Username: %s", username)

	if strings.TrimSpace(username) == "" {
		return models.User{}, errors.New("username cannot be empty")
	}
	if !role.IsValid() {
		return models.User{}, errors.New("invalid role")
	}

	existingUser, found, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		log.Printf("Repository error checking existing user: %v", err)
		return models.User{}, err
	}
	if !found {
		log.Printf("User not found for update with Username: %s", username)
		return models.User{}, errors.New("user not found")
	}

	updatedUser := models.User{
		ID:       existingUser.ID,
		Username: username,
		Role:     role,
	}

	if strings.TrimSpace(password) != "" {
		updatedUser.Password = password
	} else {
		updatedUser.Password = existingUser.Password
	}

	idString := existingUser.ID.Hex()
	if err := s.userRepo.UpdateByID(ctx, idString, updatedUser); err != nil {
		log.Printf("Repository error in UpdateUser: %v", err)
		return models.User{}, err
	}

	updatedUser.Password = ""

	log.Printf("Successfully updated user: %s", username)
	return updatedUser, nil
}

func (s *userService) DeleteUser(ctx context.Context, username string) error {
	log.Printf("Service DeleteUser called for Username: %s", username)

	if strings.TrimSpace(username) == "" {
		return errors.New("username cannot be empty")
	}

	existingUser, found, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		log.Printf("Repository error checking existing user: %v", err)
		return err
	}
	if !found {
		log.Printf("User not found for delete with Username: %s", username)
		return errors.New("user not found")
	}

	idString := existingUser.ID.Hex()
	if err := s.userRepo.DeleteByID(ctx, idString); err != nil {
		log.Printf("Repository error in DeleteUser: %v", err)
		return err
	}

	log.Printf("Successfully deleted user: %s", username)
	return nil
}
