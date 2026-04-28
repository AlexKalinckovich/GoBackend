package user

import (
	"strings"
	"sync"

	"github.com/brota/gobackend/internal/custom_errors/validation"
)

type Validator struct{}

func NewUserValidator() *Validator {
	return &Validator{}
}

func (v *Validator) validateName(name string, aggErr *validation.AggregateError, mu *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	if strings.TrimSpace(name) == "" {
		mu.Lock()
		aggErr.Add("name", name, "name is required")
		mu.Unlock()
	}
}

func (v *Validator) validateSurname(surname string, aggErr *validation.AggregateError, mu *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	if strings.TrimSpace(surname) == "" {
		mu.Lock()
		aggErr.Add("surname", surname, "surname is required")
		mu.Unlock()
	}
}

func (v *Validator) validateAge(age int, aggErr *validation.AggregateError, mu *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	if age <= 0 || age > 150 {
		mu.Lock()
		aggErr.Add("age", age, "age must be between 1 and 150")
		mu.Unlock()
	}
}

func (v *Validator) ValidateCreateOrUpdate(req CreateUserRequest, errCh chan<- *validation.AggregateError) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	aggErr := validation.NewAggregateError()

	wg.Add(3)
	go v.validateName(req.Name, aggErr, &mu, &wg)
	go v.validateSurname(req.Surname, aggErr, &mu, &wg)
	go v.validateAge(req.Age, aggErr, &mu, &wg)

	wg.Wait()

	if aggErr.HasErrors() {
		errCh <- aggErr
	} else {
		errCh <- nil
	}
}

func (v *Validator) ValidatePatch(req PatchUserRequest, errCh chan<- *validation.AggregateError) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	aggErr := validation.NewAggregateError()

	if req.Name != nil {
		wg.Add(1)
		go v.validateName(*req.Name, aggErr, &mu, &wg)
	}
	if req.Surname != nil {
		wg.Add(1)
		go v.validateSurname(*req.Surname, aggErr, &mu, &wg)
	}
	if req.Age != nil {
		wg.Add(1)
		go v.validateAge(*req.Age, aggErr, &mu, &wg)
	}

	wg.Wait()

	if aggErr.HasErrors() {
		errCh <- aggErr
	} else {
		errCh <- nil
	}
}
