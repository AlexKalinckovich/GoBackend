package service

import (
	"github.com/brota/gobackend/internal/domain/user"
	validationuser "github.com/brota/gobackend/internal/validation/user"
)

type UserService struct {
	validator validationuser.Validator
}

func NewUserService(validator validationuser.Validator) *UserService {
	return &UserService{validator: validator}
}

func (s *UserService) CreateUser(entity user.Entity) error {
	err := s.validate(entity)
	return err
}

func (s *UserService) validate(entity user.Entity) error {
	return s.validator.Validate(entity)
}
