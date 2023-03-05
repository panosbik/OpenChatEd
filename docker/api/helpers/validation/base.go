package validation

import "github.com/go-playground/validator/v10"

type ErrorResponse struct {
	FailedField string `json:"failed_field"`
	Tag         string `json:"tag"`
	Value       any    `json:"value"`
}

func NewValidator() *validator.Validate {
	validator := validator.New()
	// Register here Custom Validation
	return validator
}

var validation *validator.Validate = NewValidator()

func ValidateStruct(inObj any) (errors []*ErrorResponse) {
	err := validation.Struct(inObj)

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = err.Field()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}
