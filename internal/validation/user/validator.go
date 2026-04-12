package user

import (
	validation "github.com/brota/gobackend/internal/custom_errors/validation"
	"strings"
	"unicode"

	domainuser "github.com/brota/gobackend/internal/domain/user"
)

type Validator struct{}

func NewValidator() *Validator {
	return &Validator{}
}

func (v *Validator) Validate(entity domainuser.Entity) error {
	agg := v.initializeAggregate()
	agg = v.validateFirstName(agg, entity.FirstName)
	agg = v.validateLastName(agg, entity.LastName)
	agg = v.validateEmail(agg, entity.Email)
	agg = v.validatePhone(agg, entity.Phone)
	agg = v.validateAge(agg, entity.Age)
	return v.resolveValidation(agg)
}

func (v *Validator) initializeAggregate() *validation.ValidationAggregateError {
	return validation.NewValidationAggregateError()
}

func (v *Validator) validateFirstName(agg *validation.ValidationAggregateError, value string) *validation.ValidationAggregateError {
	v.checkLength(agg, "firstName", value)
	return agg
}

func (v *Validator) validateLastName(agg *validation.ValidationAggregateError, value string) *validation.ValidationAggregateError {
	v.checkLength(agg, "lastName", value)
	return agg
}

func (v *Validator) checkLength(agg *validation.ValidationAggregateError, field string, value string) {
	isTooShort := len(value) < 2
	v.appendError(agg, field, "must be at least 2 characters", isTooShort)
}

func (v *Validator) validateEmail(agg *validation.ValidationAggregateError, value string) *validation.ValidationAggregateError {
	v.checkEmailFormat(agg, "email", value)
	return agg
}

func (v *Validator) checkEmailFormat(agg *validation.ValidationAggregateError, field string, value string) {
	missingAt := !strings.Contains(value, "@")
	v.appendError(agg, field, "invalid email format", missingAt)
}

func (v *Validator) validatePhone(agg *validation.ValidationAggregateError, value string) *validation.ValidationAggregateError {
	v.checkPhoneFormat(agg, "phone", value)
	return agg
}

func (v *Validator) checkPhoneFormat(agg *validation.ValidationAggregateError, field string, value string) {
	isInvalid := v.isPhoneFormatInvalid(value)
	v.appendError(agg, field, "must contain 10 to 15 digits", isInvalid)
}

func (v *Validator) isPhoneFormatInvalid(value string) bool {
	hasInvalidLength := v.hasInvalidPhoneLength(value)
	hasNonDigits := v.containsNonDigits(value)
	return hasInvalidLength || hasNonDigits
}

func (v *Validator) hasInvalidPhoneLength(value string) bool {
	return len(value) < 10 || len(value) > 15
}

func (v *Validator) containsNonDigits(value string) bool {
	return v.scanValueForNonDigits(value)
}

func (v *Validator) scanValueForNonDigits(value string) bool {
	for _, char := range value {
		return v.isNotDigit(char)
	}
	return false
}

func (v *Validator) isNotDigit(char rune) bool {
	return !unicode.IsDigit(char)
}

func (v *Validator) validateAge(agg *validation.ValidationAggregateError, value int) *validation.ValidationAggregateError {
	v.checkAgeRange(agg, "age", value)
	return agg
}

func (v *Validator) checkAgeRange(agg *validation.ValidationAggregateError, field string, value int) {
	isOutOfRange := value <= 0 || value > 120
	v.appendError(agg, field, "must be between 1 and 120", isOutOfRange)
}

func (v *Validator) appendError(agg *validation.ValidationAggregateError, field string, message string, condition bool) {
	if condition {
		agg.AddField(field, message)
	}
}

func (v *Validator) resolveValidation(agg *validation.ValidationAggregateError) error {
	hasErrors := agg.HasErrors()
	return v.returnAggregateOrNil(agg, hasErrors)
}

func (v *Validator) returnAggregateOrNil(agg *validation.ValidationAggregateError, hasErrors bool) error {
	if hasErrors {
		return agg
	}
	return nil
}
